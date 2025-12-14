db = db.getSiblingDB('testdb');

// Reviews collection - Customer reviews for products
db.reviews.insertMany([
  // Reviews from 2023 orders
  {
    customer_id: 1,
    product_id: 1,
    order_year: 2023,
    rating: 5,
    title: "Excellent laptop!",
    comment: "Best laptop I've ever owned. Fast, reliable, and great battery life.",
    verified_purchase: true,
    review_date: new Date("2023-11-20"),
    helpful_count: 15
  },
  {
    customer_id: 2,
    product_id: 3,
    order_year: 2023,
    rating: 4,
    title: "Solid desk",
    comment: "Good quality desk, easy to assemble. A bit heavy though.",
    verified_purchase: true,
    review_date: new Date("2023-10-25"),
    helpful_count: 8
  },
  {
    customer_id: 3,
    product_id: 4,
    order_year: 2023,
    rating: 5,
    title: "Perfect gaming chair",
    comment: "Super comfortable for long gaming sessions. Highly recommend!",
    verified_purchase: true,
    review_date: new Date("2023-09-25"),
    helpful_count: 22
  },
  {
    customer_id: 4,
    product_id: 6,
    order_year: 2023,
    rating: 5,
    title: "Beautiful monitor",
    comment: "Colors are vibrant, perfect for photo editing.",
    verified_purchase: true,
    review_date: new Date("2023-12-28"),
    helpful_count: 12
  },
  // Reviews from 2024 orders
  {
    customer_id: 1,
    product_id: 6,
    order_year: 2024,
    rating: 5,
    title: "Great second monitor",
    comment: "Bought a second one for dual monitor setup. Perfect match!",
    verified_purchase: true,
    review_date: new Date("2024-11-18"),
    helpful_count: 7
  },
  {
    customer_id: 2,
    product_id: 7,
    order_year: 2024,
    rating: 4,
    title: "Good mechanical keyboard",
    comment: "Nice tactile feedback. A bit loud for office use.",
    verified_purchase: true,
    review_date: new Date("2024-10-22"),
    helpful_count: 9
  },
  {
    customer_id: 4,
    product_id: 1,
    order_year: 2024,
    rating: 5,
    title: "Upgraded from old laptop",
    comment: "Amazing performance upgrade. Worth every penny!",
    verified_purchase: true,
    review_date: new Date("2024-12-29"),
    helpful_count: 18
  },
  {
    customer_id: 5,
    product_id: 3,
    order_year: 2024,
    rating: 4,
    title: "Nice desk for home office",
    comment: "Perfect size for my home office setup.",
    verified_purchase: true,
    review_date: new Date("2024-12-01"),
    helpful_count: 5
  },
  {
    customer_id: 6,
    product_id: 1,
    order_year: 2024,
    rating: 5,
    title: "Best purchase of the year",
    comment: "Bought two for my business. Employees love them!",
    verified_purchase: true,
    review_date: new Date("2024-12-10"),
    helpful_count: 14
  },
  {
    customer_id: 7,
    product_id: 4,
    order_year: 2024,
    rating: 5,
    title: "So comfortable!",
    comment: "My back pain is gone after switching to this chair.",
    verified_purchase: true,
    review_date: new Date("2024-11-20"),
    helpful_count: 11
  }
]);

// Customer Interactions collection - Marketing/support data
db.customer_interactions.insertMany([
  {
    customer_id: 1,
    interaction_type: "support_ticket",
    subject: "Question about warranty",
    status: "resolved",
    created_at: new Date("2024-11-10"),
    resolved_at: new Date("2024-11-11")
  },
  {
    customer_id: 1,
    interaction_type: "newsletter_subscription",
    status: "active",
    created_at: new Date("2023-11-15")
  },
  {
    customer_id: 2,
    interaction_type: "newsletter_subscription",
    status: "active",
    created_at: new Date("2023-10-18")
  },
  {
    customer_id: 3,
    interaction_type: "support_ticket",
    subject: "Delivery question",
    status: "resolved",
    created_at: new Date("2023-09-18"),
    resolved_at: new Date("2023-09-19")
  },
  {
    customer_id: 4,
    interaction_type: "phone_call",
    subject: "Product inquiry",
    status: "completed",
    created_at: new Date("2024-12-20")
  },
  {
    customer_id: 6,
    interaction_type: "support_ticket",
    subject: "Bulk order discount",
    status: "resolved",
    created_at: new Date("2024-12-01"),
    resolved_at: new Date("2024-12-02")
  }
]);

print("MongoDB initialized with:");
print("- " + db.reviews.countDocuments() + " reviews");
print("- " + db.customer_interactions.countDocuments() + " customer interactions");
