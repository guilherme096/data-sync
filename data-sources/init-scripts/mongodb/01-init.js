db = db.getSiblingDB('testdb');

// Reviews collection - linked to clients and orders across PostgreSQL and MySQL
db.reviews.insertMany([
  // Reviews for US region orders (PostgreSQL)
  {
    client_id: 1,
    order_id: 1,
    region: "US",
    rating: 5,
    comment: "Excellent service! Product arrived quickly and works perfectly.",
    verified_purchase: true,
    review_date: new Date("2024-11-16")
  },
  {
    client_id: 1,
    order_id: 2,
    region: "US",
    rating: 4,
    comment: "Good quality, happy with the purchase.",
    verified_purchase: true,
    review_date: new Date("2024-11-29")
  },
  {
    client_id: 2,
    order_id: 3,
    region: "US",
    rating: 5,
    comment: "Amazing product! Highly recommend to everyone.",
    verified_purchase: true,
    review_date: new Date("2024-11-19")
  },
  {
    client_id: 3,
    order_id: 5,
    region: "US",
    rating: 4,
    comment: "Very satisfied with my order. Will buy again!",
    verified_purchase: true,
    review_date: new Date("2024-11-21")
  },
  {
    client_id: 4,
    order_id: 6,
    region: "US",
    rating: 5,
    comment: "Perfect! Exactly what I needed.",
    verified_purchase: true,
    review_date: new Date("2024-11-26")
  },
  // Reviews for EU region orders (MySQL)
  {
    client_id: 1,
    order_id: 1,
    region: "EU",
    rating: 5,
    comment: "Excelente produto! Muito satisfeito com a compra.",
    verified_purchase: true,
    review_date: new Date("2024-11-17")
  },
  {
    client_id: 1,
    order_id: 2,
    region: "EU",
    rating: 4,
    comment: "Bom produto, entrega rápida.",
    verified_purchase: true,
    review_date: new Date("2024-11-30")
  },
  {
    client_id: 2,
    order_id: 3,
    region: "EU",
    rating: 5,
    comment: "Perfeito! Recomendo a todos.",
    verified_purchase: true,
    review_date: new Date("2024-11-20")
  },
  {
    client_id: 3,
    order_id: 4,
    region: "EU",
    rating: 4,
    comment: "Sehr gut! Schnelle Lieferung.",
    verified_purchase: true,
    review_date: new Date("2024-11-22")
  },
  {
    client_id: 4,
    order_id: 6,
    region: "EU",
    rating: 5,
    comment: "Excellent produit! Je suis très content.",
    verified_purchase: true,
    review_date: new Date("2024-11-27")
  },
  {
    client_id: 5,
    order_id: 7,
    region: "EU",
    rating: 3,
    comment: "Prodotto buono ma la spedizione è stata lenta.",
    verified_purchase: true,
    review_date: new Date("2024-11-28")
  },
  {
    client_id: 6,
    order_id: 9,
    region: "EU",
    rating: 4,
    comment: "Dobry produkt, polecam!",
    verified_purchase: true,
    review_date: new Date("2024-11-24")
  }
]);

print("MongoDB initialized with " + db.reviews.countDocuments() + " reviews");
