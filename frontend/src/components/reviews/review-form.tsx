'use client';

import React from 'react';
import { useForm } from 'react-hook-form';
import { Star, Upload, X } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { RatingStars } from './rating-stars';
import { CreateReviewRequest, UpdateReviewRequest, Review } from '@/services/review';
import { cn } from '@/lib/utils';

interface ReviewFormData {
  rating: number;
  title: string;
  comment: string;
  images?: string[];
}

interface ReviewFormProps {
  productId: string;
  orderId?: string;
  existingReview?: Review;
  onSubmit: (data: CreateReviewRequest | UpdateReviewRequest) => Promise<void>;
  onCancel?: () => void;
  isLoading?: boolean;
  className?: string;
}

export function ReviewForm({
  productId,
  orderId,
  existingReview,
  onSubmit,
  onCancel,
  isLoading = false,
  className,
}: ReviewFormProps) {
  const [selectedImages, setSelectedImages] = React.useState<string[]>(
    existingReview?.images || []
  );
  const [uploadingImages, setUploadingImages] = React.useState(false);

  const {
    register,
    handleSubmit,
    watch,
    setValue,
    formState: { errors, isValid },
  } = useForm<ReviewFormData>({
    defaultValues: {
      rating: existingReview?.rating || 0,
      title: existingReview?.title || '',
      comment: existingReview?.comment || '',
      images: existingReview?.images || [],
    },
  });

  const watchedRating = watch('rating');

  const handleRatingChange = (rating: number) => {
    setValue('rating', rating, { shouldValidate: true });
  };

  const handleImageUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files;
    if (!files) return;

    setUploadingImages(true);
    try {
      // TODO: Implement actual image upload to your storage service
      // For now, we'll create mock URLs
      const newImageUrls = Array.from(files).map((file, index) => 
        URL.createObjectURL(file) // This is just for demo, replace with actual upload
      );
      
      setSelectedImages(prev => [...prev, ...newImageUrls]);
      setValue('images', [...selectedImages, ...newImageUrls]);
    } catch (error) {
      console.error('Failed to upload images:', error);
    } finally {
      setUploadingImages(false);
    }
  };

  const removeImage = (index: number) => {
    const newImages = selectedImages.filter((_, i) => i !== index);
    setSelectedImages(newImages);
    setValue('images', newImages);
  };

  const onFormSubmit = async (data: ReviewFormData) => {
    try {
      if (existingReview) {
        // Update existing review
        await onSubmit({
          rating: data.rating,
          title: data.title,
          comment: data.comment,
          images: selectedImages,
        });
      } else {
        // Create new review
        await onSubmit({
          product_id: productId,
          order_id: orderId,
          rating: data.rating,
          title: data.title,
          comment: data.comment,
          images: selectedImages,
        });
      }
    } catch (error) {
      console.error('Failed to submit review:', error);
    }
  };

  const getRatingLabel = (rating: number) => {
    const labels = {
      1: 'Poor',
      2: 'Fair',
      3: 'Good',
      4: 'Very Good',
      5: 'Excellent',
    };
    return labels[rating as keyof typeof labels] || '';
  };

  return (
    <Card className={cn('w-full max-w-2xl', className)}>
      <CardHeader>
        <CardTitle>
          {existingReview ? 'Edit Your Review' : 'Write a Review'}
        </CardTitle>
      </CardHeader>
      
      <CardContent>
        <form onSubmit={handleSubmit(onFormSubmit)} className="space-y-6">
          {/* Rating */}
          <div className="space-y-2">
            <Label className="text-base font-semibold">Overall Rating *</Label>
            <div className="flex items-center gap-4">
              <RatingStars
                rating={watchedRating}
                size="lg"
                interactive
                onRatingChange={handleRatingChange}
              />
              {watchedRating > 0 && (
                <span className="text-sm font-medium text-gray-600">
                  {getRatingLabel(watchedRating)}
                </span>
              )}
            </div>
            {errors.rating && (
              <p className="text-sm text-red-600">Please select a rating</p>
            )}
          </div>

          {/* Title */}
          <div className="space-y-2">
            <Label htmlFor="title" className="text-base font-semibold">
              Review Title <span className="text-gray-400 text-sm">(Optional)</span>
            </Label>
            <Input
              id="title"
              placeholder="Summarize your experience (optional)"
              {...register('title', {
                maxLength: { value: 100, message: 'Title must be less than 100 characters' }
              })}
              className={errors.title ? 'border-red-500' : ''}
            />
            {errors.title && (
              <p className="text-sm text-red-600">{errors.title.message}</p>
            )}
            <p className="text-xs text-gray-500">
              If left empty, we'll generate a title based on your rating
            </p>
          </div>

          {/* Comment */}
          <div className="space-y-2">
            <Label htmlFor="comment" className="text-base font-semibold">
              Your Review <span className="text-gray-400 text-sm">(Optional)</span>
            </Label>
            <Textarea
              id="comment"
              placeholder="Tell others about your experience with this product (optional)..."
              rows={5}
              {...register('comment', {
                maxLength: { value: 1000, message: 'Review must be less than 1000 characters' }
              })}
              className={errors.comment ? 'border-red-500' : ''}
            />
            {errors.comment && (
              <p className="text-sm text-red-600">{errors.comment.message}</p>
            )}
            <p className="text-xs text-gray-500">
              You can just rate the product without writing a comment
            </p>
          </div>

          {/* Images */}
          <div className="space-y-2">
            <Label className="text-base font-semibold">Add Photos (Optional)</Label>
            <div className="space-y-4">
              {/* Image Upload */}
              <div className="flex items-center gap-4">
                <label className="cursor-pointer">
                  <input
                    type="file"
                    multiple
                    accept="image/*"
                    onChange={handleImageUpload}
                    className="hidden"
                    disabled={uploadingImages || selectedImages.length >= 5}
                  />
                  <Button
                    type="button"
                    variant="outline"
                    disabled={uploadingImages || selectedImages.length >= 5}
                    className="flex items-center gap-2"
                  >
                    <Upload className="w-4 h-4" />
                    {uploadingImages ? 'Uploading...' : 'Add Photos'}
                  </Button>
                </label>
                <span className="text-sm text-gray-500">
                  {selectedImages.length}/5 photos
                </span>
              </div>

              {/* Selected Images */}
              {selectedImages.length > 0 && (
                <div className="grid grid-cols-5 gap-2">
                  {selectedImages.map((image, index) => (
                    <div key={index} className="relative group">
                      <img
                        src={image}
                        alt={`Review image ${index + 1}`}
                        className="w-full h-20 object-cover rounded-lg border"
                      />
                      <button
                        type="button"
                        onClick={() => removeImage(index)}
                        className="absolute -top-2 -right-2 w-6 h-6 bg-red-500 text-white rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity"
                      >
                        <X className="w-4 h-4" />
                      </button>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>

          {/* Actions */}
          <div className="flex items-center gap-4 pt-4">
            <Button
              type="submit"
              disabled={isLoading || watchedRating === 0}
              className="flex-1"
            >
              {isLoading ? 'Submitting...' : existingReview ? 'Update Review' : 'Submit Review'}
            </Button>
            
            {onCancel && (
              <Button type="button" variant="outline" onClick={onCancel}>
                Cancel
              </Button>
            )}
          </div>
        </form>
      </CardContent>
    </Card>
  );
}
