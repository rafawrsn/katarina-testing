export interface ProductVariant {
  id: string;
  product_id?: string;
  sku: string;
  color: string;
  size: string;
  price: string | number;
  stock: number;
  weight?: number;
  status?: string;
}

export interface ProductImage {
  id?: string;
  image_url: string;
  is_primary?: boolean;
  sort_order?: number;
}

export interface Category {
  id: string;
  name: string;
  slug: string;
  description?: string;
  image_url?: string;
}

export interface Product {
  id: string | number;
  name: string;
  slug?: string;
  description: string;
  price: string;
  rawPrice?: number;
  image: string;
  category: string;
  category_id?: string;
  shopeeLink?: string;
  variants?: ProductVariant[];
  images?: ProductImage[];
  min_price?: string | number;
  max_price?: string | number;
  total_stock?: number;
  status?: string;
}

export const fallbackProducts: Product[] = [
  {
    id: '1',
    name: 'Katarina Embroidery Shirt',
    slug: 'katarina-embroidery-shirt',
    description: 'Minimalist shirt with special embroidery details. Self-manufactured with breathable soft cotton.',
    price: 'Rp129.000',
    rawPrice: 129000,
    image:
      'https://images.pexels.com/photos/22441297/pexels-photo-22441297.jpeg?auto=compress&cs=tinysrgb&h=900',
    category: 'Tops',
    shopeeLink: 'https://shopee.co.id/',
    variants: [
      { id: 'v1-1', sku: 'KD-SHIRT-CRM-S', color: 'Cream', size: 'S', price: 129000, stock: 25 },
      { id: 'v1-2', sku: 'KD-SHIRT-CRM-M', color: 'Cream', size: 'M', price: 129000, stock: 30 },
      { id: 'v1-3', sku: 'KD-SHIRT-CRM-L', color: 'Cream', size: 'L', price: 129000, stock: 20 },
      { id: 'v1-4', sku: 'KD-SHIRT-BGE-S', color: 'Beige', size: 'S', price: 129000, stock: 15 },
      { id: 'v1-5', sku: 'KD-SHIRT-BGE-M', color: 'Beige', size: 'M', price: 129000, stock: 20 },
    ],
  },
  {
    id: '2',
    name: 'Daily Minimalist Blouse',
    slug: 'daily-minimalist-blouse',
    description: 'Elegant Korean-inspired blouse for daily style with delicate cuff stitching.',
    price: 'Rp149.000',
    rawPrice: 149000,
    image:
      'https://images.pexels.com/photos/7945547/pexels-photo-7945547.jpeg?auto=compress&cs=tinysrgb&h=900',
    category: 'Tops',
    shopeeLink: 'https://shopee.co.id/',
    variants: [
      { id: 'v2-1', sku: 'KD-BLOUSE-WHT-S', color: 'Off White', size: 'S', price: 149000, stock: 30 },
      { id: 'v2-2', sku: 'KD-BLOUSE-WHT-M', color: 'Off White', size: 'M', price: 149000, stock: 25 },
      { id: 'v2-3', sku: 'KD-BLOUSE-BLU-S', color: 'Soft Blue', size: 'S', price: 149000, stock: 15 },
      { id: 'v2-4', sku: 'KD-BLOUSE-BLU-M', color: 'Soft Blue', size: 'M', price: 149000, stock: 18 },
    ],
  },
  {
    id: '3',
    name: 'Soft Cotton Outer',
    slug: 'soft-cotton-outer',
    description: 'Comfortable minimalist outerwear with an oversized drape and subtle lapel detail.',
    price: 'Rp179.000',
    rawPrice: 179000,
    image:
      'https://images.pexels.com/photos/9603628/pexels-photo-9603628.jpeg?auto=compress&cs=tinysrgb&h=900',
    category: 'Outerwear',
    shopeeLink: 'https://shopee.co.id/',
    variants: [
      { id: 'v3-1', sku: 'KD-OUTER-OAT-ALL', color: 'Oatmeal', size: 'All Size', price: 179000, stock: 25 },
      { id: 'v3-2', sku: 'KD-OUTER-CHR-ALL', color: 'Charcoal', size: 'All Size', price: 179000, stock: 15 },
    ],
  },
  {
    id: '4',
    name: 'Neutral Knit Cardigan',
    slug: 'neutral-knit-cardigan',
    description: 'Soft knit cardigan in warm neutral tones with custom horn-finish buttons.',
    price: 'Rp159.000',
    rawPrice: 159000,
    image:
      'https://images.pexels.com/photos/4641824/pexels-photo-4641824.jpeg?auto=compress&cs=tinysrgb&h=900',
    category: 'Outerwear',
    shopeeLink: 'https://shopee.co.id/',
    variants: [
      { id: 'v4-1', sku: 'KD-CARD-ALM-SM', color: 'Almond', size: 'S-M', price: 159000, stock: 20 },
      { id: 'v4-2', sku: 'KD-CARD-ALM-LXL', color: 'Almond', size: 'L-XL', price: 159000, stock: 18 },
      { id: 'v4-3', sku: 'KD-CARD-MCH-SM', color: 'Mocha', size: 'S-M', price: 159000, stock: 12 },
    ],
  },
  {
    id: '5',
    name: 'Linen Everyday Shirt',
    slug: 'linen-everyday-shirt',
    description: 'Breathable linen shirt for effortless daily wear, featuring fine collar embroidery.',
    price: 'Rp139.000',
    rawPrice: 139000,
    image:
      'https://images.pexels.com/photos/6373455/pexels-photo-6373455.jpeg?auto=compress&cs=tinysrgb&h=900',
    category: 'Tops',
    shopeeLink: 'https://shopee.co.id/',
    variants: [
      { id: 'v5-1', sku: 'KD-LINEN-NAT-S', color: 'Natural', size: 'S', price: 139000, stock: 20 },
      { id: 'v5-2', sku: 'KD-LINEN-NAT-M', color: 'Natural', size: 'M', price: 139000, stock: 25 },
      { id: 'v5-3', sku: 'KD-LINEN-OLV-S', color: 'Olive', size: 'S', price: 139000, stock: 14 },
    ],
  },
  {
    id: '6',
    name: 'Beige Tailored Coat',
    slug: 'beige-tailored-coat',
    description: 'Refined tailored coat with a calm minimalist silhouette and inner pocket embroidery.',
    price: 'Rp199.000',
    rawPrice: 199000,
    image:
      'https://images.pexels.com/photos/4456728/pexels-photo-4456728.jpeg?auto=compress&cs=tinysrgb&h=900',
    category: 'Outerwear',
    shopeeLink: 'https://shopee.co.id/',
    variants: [
      { id: 'v6-1', sku: 'KD-COAT-BGE-M', color: 'Warm Beige', size: 'M', price: 199000, stock: 12 },
      { id: 'v6-2', sku: 'KD-COAT-BGE-L', color: 'Warm Beige', size: 'L', price: 199000, stock: 10 },
      { id: 'v6-3', sku: 'KD-COAT-BLK-M', color: 'Black', size: 'M', price: 199000, stock: 15 },
    ],
  },
];

export const products = fallbackProducts;
export const categories = ['All', 'Tops', 'Outerwear', 'Dresses', 'Bottoms'];
