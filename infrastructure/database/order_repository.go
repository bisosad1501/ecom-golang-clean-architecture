	return orders, err
}

// CountSearch counts orders based on search criteria
func (r *orderRepository) CountSearch(ctx context.Context, params repositories.OrderSearchParams) (int64, error) {
	query := r.db.WithContext(ctx).Model(&entities.Order{})

	// Apply the same filters as Search method
	if params.UserID != nil {
		query = query.Where("user_id = ?", *params.UserID)
	}

	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}

	if params.PaymentStatus != nil {
		query = query.Where("payment_status = ?", *params.PaymentStatus)
	}

	if params.StartDate != nil {
		query = query.Where("created_at >= ?", *params.StartDate)
	}

	if params.EndDate != nil {
		query = query.Where("created_at <= ?", *params.EndDate)
	}

	if params.MinTotal != nil {
		query = query.Where("total >= ?", *params.MinTotal)
	}

	if params.MaxTotal != nil {
		query = query.Where("total <= ?", *params.MaxTotal)
	}

	var count int64
	err := query.Count(&count).Error
	return count, err
}