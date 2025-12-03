db = db.getSiblingDB('testdb');

db.reviews.insertMany([
  {
    product_id: 1,
    customer_name: "João Silva",
    rating: 5,
    comment: "Excellent laptop! Very fast and reliable.",
    verified_purchase: true,
    review_date: new Date("2024-11-15")
  },
  {
    product_id: 2,
    customer_name: "Maria Santos",
    rating: 4,
    comment: "Good mouse, comfortable grip.",
    verified_purchase: true,
    review_date: new Date("2024-11-20")
  },
  {
    product_id: 3,
    customer_name: "Pedro Costa",
    rating: 5,
    comment: "Best mechanical keyboard I've owned!",
    verified_purchase: true,
    review_date: new Date("2024-11-18")
  },
  {
    product_id: 4,
    customer_name: "Ana Rodrigues",
    rating: 4,
    comment: "Great monitor, colors are vibrant.",
    verified_purchase: false,
    review_date: new Date("2024-11-22")
  },
  {
    product_id: 1,
    customer_name: "Carlos Ferreira",
    rating: 5,
    comment: "Worth every penny, highly recommend!",
    verified_purchase: true,
    review_date: new Date("2024-11-25")
  },
  {
    product_id: 6,
    customer_name: "Luisa Alves",
    rating: 3,
    comment: "Webcam is okay, image quality could be better.",
    verified_purchase: true,
    review_date: new Date("2024-11-28")
  }
]);

print("MongoDB initialized with " + db.reviews.countDocuments() + " reviews");
