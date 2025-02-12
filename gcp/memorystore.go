package gcp

import (
	"context"
	"errors"
	"fmt"
	"iter"

	"github.com/jarrodhroberson/ossgo/functions/must"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"

	redis "cloud.google.com/go/redis/apiv1"
	"cloud.google.com/go/redis/apiv1/redispb"
)

func newCloudRedisClient(ctx context.Context) (*redis.CloudRedisClient, error) {
	client, err := redis.NewCloudRedisClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Redis client: %v", err)
	}
	return client, nil
}

func getRedisInstance(ctx context.Context, name string) (*redispb.Instance, error) {
	c, err := newCloudRedisClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Redis client: %v", err)
	}
	defer c.Close()

	req := &redispb.GetInstanceRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/instances/%s", must.Must(ProjectId()), must.Must(Region()), name),
	}

	resp, err := c.GetInstance(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get Redis instance: %v", err)
	}

	return resp, nil
}

func createRedisInstance(ctx context.Context, name string, tier redispb.Instance_Tier, memorySizeGb int32, location string) {
	c, err := newCloudRedisClient(ctx)
	if err != nil {
		panic(err)
	}
	defer func(c *redis.CloudRedisClient) {
		err := c.Close()
		if err != nil {
			log.Warn().Err(err).Msg("failed to close Redis client")
		}
	}(c)

	req := &redispb.CreateInstanceRequest{
		Parent:     fmt.Sprintf("projects/%s/locations/%s", must.Must(ProjectId()), must.Must(Region())),
		InstanceId: name,
		Instance: &redispb.Instance{
			Tier:         tier,
			MemorySizeGb: memorySizeGb,
		},
	}
	op, err := c.CreateInstance(ctx, req)
	if err != nil {
		panic(err)
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		panic(err)
	}

	log.Info().Msgf("Redis Instance Endpoing: %s", resp.GetReadEndpoint())
}

func deleteRedisInstance(ctx context.Context, name string) {
	c, err := newCloudRedisClient(ctx)
	if err != nil {
		panic(err)
	}
	defer func(c *redis.CloudRedisClient) {
		err := c.Close()
		if err != nil {
			log.Warn().Err(err).Msg("failed to close Redis client")
		}
	}(c)

	req := &redispb.DeleteInstanceRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/instances/%s", must.Must(ProjectId()), must.Must(Region()), name),
	}

	op, err := c.DeleteInstance(ctx, req)
	if err != nil {
		panic(err)
	}

	err = op.Wait(ctx)
	if err != nil {
		panic(err)
	}

	log.Info().Msgf("Redis Instance %s deleted successfully", name)
}

func listRedisInstances(ctx context.Context) (iter.Seq[string], error) {
	c, err := newCloudRedisClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Redis client: %v", err)
	}
	defer func(c *redis.CloudRedisClient) {
		err := c.Close()
		if err != nil {
			log.Warn().Err(err).Msg("failed to close Redis client")
		}
	}(c)

	req := &redispb.ListInstancesRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", must.Must(ProjectId()), must.Must(Region())),
	}

	return func(yield func(string) bool) {
		iiter := c.ListInstances(ctx, req)
		for {
			instance, err := iiter.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				log.Error().Err(err).Msg("Error iterating through Redis instances")
				return
			}
			if !yield(instance.GetName()) {
				break
			}
		}
	}, nil
}
