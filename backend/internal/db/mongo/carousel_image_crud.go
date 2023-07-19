package mongo

import (
	"bar/internal/models"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (b *Backend) CreateCarouselImage(ci *models.CarouselImage) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	_, err := b.db.Collection(CarouselImagesCollection).InsertOne(ctx, ci)
	if err != nil {
		return err
	}

	return nil
}

func (b *Backend) GetCarouselImage(id string) (*models.CarouselImage, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	var ci models.CarouselImage
	err := b.db.Collection(CarouselImagesCollection).FindOne(ctx,
		bson.M{
			"id": uuid.MustParse(id),

			"$or": []bson.M{
				{
					"deleted_at": bson.M{
						"$exists": false,
					},
				},
				{
					"deleted_at": nil,
				},
			},
		},
	).Decode(&ci)
	if err != nil {
		return nil, err
	}

	return &ci, nil
}

func (b *Backend) UpdateCarouselImage(ci *models.CarouselImage) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(CarouselImagesCollection).FindOneAndUpdate(ctx,
		bson.M{
			"id": ci.Id,

			"$or": []bson.M{
				{
					"deleted_at": bson.M{
						"$exists": false,
					},
				},
				{
					"deleted_at": nil,
				},
			},
		},
		bson.M{
			"$set": ci,
		},
		options.FindOneAndUpdate().SetUpsert(true))
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) MarkDeleteCarouselImage(id, by string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(CarouselImagesCollection).FindOneAndUpdate(ctx,
		bson.M{
			"id": uuid.MustParse(id),
		},
		bson.M{
			"$set": bson.M{
				"deleted_at": time.Now().Unix(),
				"deleted_by": uuid.MustParse(by),
			},
		},
		options.FindOneAndUpdate().SetUpsert(false))
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) UnMarkDeleteCarouselImage(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(CarouselImagesCollection).FindOneAndUpdate(ctx,
		bson.M{
			"id": uuid.MustParse(id),
		},
		bson.M{
			"$set": bson.M{
				"deleted_at": nil,
				"deleted_by": nil,
			},
		},
		options.FindOneAndUpdate().SetUpsert(false))
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) DeleteCarouselImage(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(CarouselImagesCollection).FindOneAndDelete(ctx,
		bson.M{
			"id": uuid.MustParse(id),
		},
	)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) RestoreCarouselImage(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(CarouselImagesCollection).FindOneAndUpdate(ctx,
		bson.M{
			"id": uuid.MustParse(id),
		},
		bson.M{
			"$unset": bson.M{
				"deleted_at": "",
				"deleted_by": "",
			},
		},
	)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) GetDeletedCarouselImages(page uint64, size uint64) ([]*models.CarouselImage, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	var accs []*models.CarouselImage
	cursor, err := b.db.Collection(CarouselImagesCollection).Find(ctx,
		bson.M{
			"deleted_at": bson.M{
				"$ne": nil,
			},
		},
		options.Find().SetSkip(int64(page*size)).SetLimit(int64(size)))
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &accs); err != nil {
		return nil, err
	}

	return accs, nil
}

func (b *Backend) CountDeletedCarouselImages() (int64, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	count, err := b.db.Collection(CarouselImagesCollection).CountDocuments(ctx, bson.M{
		"deleted_at": bson.M{
			"$ne": nil,
		},
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}
