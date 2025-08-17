SELECT p.code, p.price, c.name
FROM product AS p
INNER JOIN category as c
  ON p.category_id = c.id
WHERE p.price < ? AND c.id IN ?;
