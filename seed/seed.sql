INSERT INTO products (name, type, price)
VALUES
  ('Bayam', 'Sayuran', 12000),
  ('Wortel', 'Sayuran', 15000),
  ('Ayam Fillet', 'Protein', 45000),
  ('Telur Ayam 1 Lusin', 'Protein', 32000),
  ('Apel Fuji', 'Buah', 28000),
  ('Pisang Cavendish', 'Buah', 22000),
  ('Keripik Kentang', 'Snack', 18000),
  ('Biskuit Coklat', 'Snack', 16000)
ON CONFLICT DO NOTHING;

