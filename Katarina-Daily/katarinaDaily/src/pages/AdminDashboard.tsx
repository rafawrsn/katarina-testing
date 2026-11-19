import React, { useEffect, useState, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { motion, AnimatePresence } from 'framer-motion';
import {
  ArrowLeft,
  Package,
  ShoppingBag,
  Plus,
  Edit2,
  Trash2,
  Search,
  CheckCircle2,
  AlertCircle,
  X,
  ExternalLink,
  Layers,
  Sparkles,
  ShieldCheck,
  RefreshCw,
  SlidersHorizontal,
  Tag,
  Boxes,
  Check,
  AlertTriangle
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { getAllOrders, updateOrderStatus, deleteOrder } from '@/services/adminOrderService';
import { getProducts, getCategoryObjects, formatIDR, type Product, type Category, type ProductVariant } from '@/services/catalogService';
import {
  createProduct,
  updateProduct,
  deleteProduct,
  addProductImage,
  addProductVariant,
  updateProductVariant,
  deleteProductVariant,
  createCategory
} from '@/services/adminProductService';

// Preset sample images for quick convenience
const PRESET_IMAGES = [
  'https://images.pexels.com/photos/22441297/pexels-photo-22441297.jpeg?auto=compress&cs=tinysrgb&h=900',
  'https://images.pexels.com/photos/7945547/pexels-photo-7945547.jpeg?auto=compress&cs=tinysrgb&h=900',
  'https://images.pexels.com/photos/9603628/pexels-photo-9603628.jpeg?auto=compress&cs=tinysrgb&h=900',
  'https://images.pexels.com/photos/8070398/pexels-photo-8070398.jpeg?auto=compress&cs=tinysrgb&h=900',
  'https://images.pexels.com/photos/6373455/pexels-photo-6373455.jpeg?auto=compress&cs=tinysrgb&h=900',
  'https://images.pexels.com/photos/4456728/pexels-photo-4456728.jpeg?auto=compress&cs=tinysrgb&h=900',
];

interface VariantFormItem {
  id?: string;
  sku: string;
  color: string;
  size: string;
  price: number | string;
  stock: number;
  weight: number;
  status: 'active' | 'inactive';
}

interface ProductFormData {
  name: string;
  slug: string;
  category_id: string;
  description: string;
  status: 'active' | 'draft' | 'inactive';
  primaryImageUrl: string;
  variants: VariantFormItem[];
}

const DEFAULT_FORM_DATA: ProductFormData = {
  name: '',
  slug: '',
  category_id: '',
  description: '',
  status: 'active',
  primaryImageUrl: 'https://images.pexels.com/photos/22441297/pexels-photo-22441297.jpeg?auto=compress&cs=tinysrgb&h=900',
  variants: [
    { sku: '', color: 'Cream', size: 'S', price: 149000, stock: 20, weight: 250, status: 'active' },
    { sku: '', color: 'Cream', size: 'M', price: 149000, stock: 25, weight: 250, status: 'active' },
  ],
};

export default function AdminDashboard() {
  const { isAdmin, isAuthenticated, loading: authLoading, openAuthModal } = useAuth();

  const [activeTab, setActiveTab] = useState<'products' | 'orders' | 'categories'>('products');

  // Products & Categories Data
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loadingProducts, setLoadingProducts] = useState(false);

  // Orders Data
  const [orders, setOrders] = useState<any[]>([]);
  const [loadingOrders, setLoadingOrders] = useState(false);

  // Filters & Search
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState<string>('all');
  const [selectedStatus, setSelectedStatus] = useState<string>('all');
  const [stockFilter, setStockFilter] = useState<'all' | 'in_stock' | 'low_stock' | 'out_of_stock'>('all');
  const [sortBy, setSortBy] = useState<'newest' | 'name_asc' | 'name_desc' | 'price_asc' | 'price_desc' | 'stock_desc'>('newest');
  const [minPrice, setMinPrice] = useState<string>('');
  const [maxPrice, setMaxPrice] = useState<string>('');
  const [minPriceError, setMinPriceError] = useState<string | null>(null);

  // Modals
  const [isFormModalOpen, setIsFormModalOpen] = useState(false);
  const [editingProduct, setEditingProduct] = useState<Product | null>(null);
  const [formData, setFormData] = useState<ProductFormData>(DEFAULT_FORM_DATA);
  const [formSubmitting, setFormSubmitting] = useState(false);
  const [formError, setFormError] = useState<string | null>(null);

  // Variant Manager Modal
  const [managingVariantProduct, setManagingVariantProduct] = useState<Product | null>(null);
  const [variantSubmitting, setVariantSubmitting] = useState(false);

  // Delete Confirmation Modal
  const [deletingProduct, setDeletingProduct] = useState<Product | null>(null);
  const [deletingOrder, setDeletingOrder] = useState<any | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);

  // New Category Modal
  const [isNewCatModalOpen, setIsNewCatModalOpen] = useState(false);
  const [newCatName, setNewCatName] = useState('');
  const [newCatDescription, setNewCatDescription] = useState('');
  const [catSubmitting, setCatSubmitting] = useState(false);

  // Toast Notification
  const [toast, setToast] = useState<{ message: string; type: 'success' | 'error' | 'info' } | null>(null);

  const showToast = (message: string, type: 'success' | 'error' | 'info' = 'success') => {
    setToast({ message, type });
    setTimeout(() => {
      setToast(null);
    }, 3500);
  };

  // Fetch initial data
  const loadCategories = async () => {
    try {
      const cats = await getCategoryObjects();
      setCategories(cats);
      if (cats.length > 0 && !formData.category_id) {
        setFormData((prev) => ({ ...prev, category_id: cats[0].id }));
      }
    } catch (err) {
      console.error('Failed to load categories', err);
    }
  };

  const loadProducts = async () => {
    setLoadingProducts(true);
    try {
      const list = await getProducts();
      setProducts(list);
    } catch (err) {
      console.error('Failed to load products', err);
      showToast('Failed to load catalog products', 'error');
    } finally {
      setLoadingProducts(false);
    }
  };

  const loadOrders = async () => {
    setLoadingOrders(true);
    try {
      const list = await getAllOrders();
      setOrders(list);
    } catch (err) {
      console.error('Failed to load orders', err);
    } finally {
      setLoadingOrders(false);
    }
  };

  useEffect(() => {
    loadCategories();
    loadProducts();
    loadOrders();
  }, []);

  useEffect(() => {
    if (activeTab === 'orders') loadOrders();
    if (activeTab === 'products') loadProducts();
    if (activeTab === 'categories') loadCategories();
  }, [activeTab]);

  // Open Create Modal
  const handleOpenCreateModal = () => {
    const defaultCatId = categories[0]?.id || '';
    setFormData({
      name: '',
      slug: '',
      category_id: defaultCatId,
      description: '',
      status: 'active',
      primaryImageUrl: PRESET_IMAGES[0],
      variants: [
        { sku: `KD-${Date.now().toString().slice(-4)}-S`, color: 'Cream', size: 'S', price: 149000, stock: 20, weight: 250, status: 'active' },
        { sku: `KD-${Date.now().toString().slice(-4)}-M`, color: 'Cream', size: 'M', price: 149000, stock: 25, weight: 250, status: 'active' },
      ],
    });
    setEditingProduct(null);
    setFormError(null);
    setIsFormModalOpen(true);
  };

  // Open Edit Modal
  const handleOpenEditModal = (product: Product) => {
    setEditingProduct(product);
    setFormError(null);

    // Match category ID
    let catId = product.category_id;
    if (!catId) {
      const matchedCat = categories.find((c) => c.name.toLowerCase() === product.category.toLowerCase());
      catId = matchedCat?.id || categories[0]?.id || '';
    }

    const mappedVariants: VariantFormItem[] = (product.variants || []).map((v) => ({
      id: v.id,
      sku: v.sku || `KD-${product.id}-${v.size}`,
      color: v.color || 'Standard',
      size: v.size || 'All Size',
      price: v.price || product.rawPrice || 149000,
      stock: v.stock !== undefined ? v.stock : 10,
      weight: v.weight || 250,
      status: (v.status as 'active' | 'inactive') || 'active',
    }));

    setFormData({
      name: product.name,
      slug: product.slug || product.name.toLowerCase().replace(/[^a-z0-9]+/g, '-'),
      category_id: catId,
      description: product.description || '',
      status: (product.status as 'active' | 'draft' | 'inactive') || 'active',
      primaryImageUrl: product.image || PRESET_IMAGES[0],
      variants: mappedVariants.length > 0 ? mappedVariants : [
        { sku: `KD-${product.id}-ALL`, color: 'Standard', size: 'All Size', price: product.rawPrice || 149000, stock: product.total_stock || 10, weight: 250, status: 'active' }
      ],
    });

    setIsFormModalOpen(true);
  };

  // Auto slug generation on name change
  const handleNameChange = (name: string) => {
    if (!editingProduct) {
      const generatedSlug = name
        .toLowerCase()
        .trim()
        .replace(/[^\w\s-]/g, '')
        .replace(/[\s_-]+/g, '-')
        .replace(/^-+|-+$/g, '');
      setFormData((prev) => ({ ...prev, name, slug: generatedSlug }));
    } else {
      setFormData((prev) => ({ ...prev, name }));
    }
  };

  // Add a variant row in form
  const handleAddVariantRow = () => {
    const baseColor = formData.variants[formData.variants.length - 1]?.color || 'Cream';
    const basePrice = formData.variants[formData.variants.length - 1]?.price || 149000;
    const newSku = `KD-${Date.now().toString().slice(-4)}-L`;

    setFormData((prev) => ({
      ...prev,
      variants: [
        ...prev.variants,
        { sku: newSku, color: baseColor, size: 'L', price: basePrice, stock: 15, weight: 250, status: 'active' },
      ],
    }));
  };

  // Remove a variant row in form
  const handleRemoveVariantRow = (index: number) => {
    if (formData.variants.length <= 1) {
      showToast('A product must contain at least one variant', 'info');
      return;
    }
    setFormData((prev) => ({
      ...prev,
      variants: prev.variants.filter((_, i) => i !== index),
    }));
  };

  // Update a variant row field
  const handleUpdateVariantRow = (index: number, field: keyof VariantFormItem, value: any) => {
    setFormData((prev) => {
      const next = [...prev.variants];
      next[index] = { ...next[index], [field]: value };
      return { ...prev, variants: next };
    });
  };

  // Quick preset sizes generator
  const handleGenerateSizes = (color: string) => {
    const sizes = ['S', 'M', 'L'];
    const basePrice = formData.variants[0]?.price || 149000;
    const newVariants: VariantFormItem[] = sizes.map((size) => ({
      sku: `KD-${formData.slug.toUpperCase().slice(0, 6) || 'ITEM'}-${color.slice(0, 3).toUpperCase()}-${size}`,
      color,
      size,
      price: basePrice,
      stock: 20,
      weight: 250,
      status: 'active',
    }));
    setFormData((prev) => ({ ...prev, variants: newVariants }));
  };

  // Save Product (Create or Update)
  const handleSaveProduct = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormError(null);

    // Validation
    if (!formData.name.trim()) {
      setFormError('Product title is required.');
      return;
    }
    if (!formData.category_id) {
      setFormError('Please select a valid category.');
      return;
    }
    if (formData.variants.length === 0) {
      setFormError('At least one variant (color/size/price) is required.');
      return;
    }

    setFormSubmitting(true);
    try {
      if (editingProduct) {
        // ---- UPDATE EXISTING PRODUCT ----
        const productId = String(editingProduct.id);

        // 1. Update basic product info
        await updateProduct(productId, {
          name: formData.name.trim(),
          slug: formData.slug.trim(),
          category_id: formData.category_id,
          description: formData.description.trim(),
          status: formData.status,
        });

        // 2. Add new primary image if changed
        if (formData.primaryImageUrl && formData.primaryImageUrl !== editingProduct.image) {
          try {
            await addProductImage(productId, {
              image_url: formData.primaryImageUrl,
              is_primary: true,
              sort_order: 0,
            });
          } catch (imgErr) {
            console.warn('Image update warning:', imgErr);
          }
        }

        // 3. Update or create variants
        for (const v of formData.variants) {
          if (v.id) {
            // Update existing variant
            try {
              await updateProductVariant(v.id, {
                sku: v.sku,
                color: v.color,
                size: v.size,
                price: Number(v.price),
                stock: Number(v.stock),
                weight: Number(v.weight),
                status: v.status,
              });
            } catch (err) {
              console.warn(`Could not update variant ${v.id}:`, err);
            }
          } else {
            // Create newly added variant
            try {
              await addProductVariant(productId, {
                sku: v.sku || `KD-${productId.slice(0, 4)}-${v.size}`,
                color: v.color,
                size: v.size,
                price: Number(v.price),
                stock: Number(v.stock),
                weight: Number(v.weight),
                status: v.status,
              });
            } catch (err) {
              console.warn('Could not add variant:', err);
            }
          }
        }

        showToast(`"${formData.name}" successfully updated!`, 'success');
      } else {
        // ---- CREATE NEW PRODUCT ----
        // 1. Create base product in backend
        const created = await createProduct({
          category_id: formData.category_id,
          name: formData.name.trim(),
          slug: formData.slug.trim() || undefined,
          description: formData.description.trim(),
          status: formData.status,
        });

        const newProductId = created.id;

        // 2. Add primary image
        if (formData.primaryImageUrl) {
          try {
            await addProductImage(newProductId, {
              image_url: formData.primaryImageUrl,
              is_primary: true,
              sort_order: 0,
            });
          } catch (imgErr) {
            console.warn('Could not register primary image:', imgErr);
          }
        }

        // 3. Add all variants
        for (const v of formData.variants) {
          try {
            await addProductVariant(newProductId, {
              sku: v.sku || `KD-${newProductId.slice(0, 4)}-${v.size}`,
              color: v.color || 'Standard',
              size: v.size || 'All Size',
              price: Number(v.price),
              stock: Number(v.stock),
              weight: Number(v.weight) || 250,
              status: v.status || 'active',
            });
          } catch (varErr) {
            console.warn('Failed to add variant during creation:', varErr);
          }
        }

        showToast(`New product "${formData.name}" created with ${formData.variants.length} variants!`, 'success');
      }

      setIsFormModalOpen(false);
      await loadProducts();
    } catch (err: any) {
      console.error('Error saving product:', err);
      setFormError(err.message || 'Failed to save product. Check backend logs.');
    } finally {
      setFormSubmitting(false);
    }
  };

  // Delete product confirmation
  const handleConfirmDeleteProduct = async () => {
    if (!deletingProduct) return;
    setIsDeleting(true);
    try {
      await deleteProduct(String(deletingProduct.id));
      showToast(`Product "${deletingProduct.name}" removed from catalog`, 'success');
      setDeletingProduct(null);
      await loadProducts();
    } catch (err: any) {
      console.error('Failed to delete product', err);
      showToast(err.message || 'Failed to delete product', 'error');
    } finally {
      setIsDeleting(false);
    }
  };

  // Quick variant stock / price update from variant manager modal
  const handleUpdateQuickVariant = async (variantId: string, updates: Partial<ProductVariant>) => {
    setVariantSubmitting(true);
    try {
      await updateProductVariant(variantId, {
        stock: updates.stock !== undefined ? Number(updates.stock) : undefined,
        price: updates.price !== undefined ? Number(updates.price) : undefined,
        status: updates.status as 'active' | 'inactive',
      });

      // Update in local modal state
      if (managingVariantProduct) {
        const nextVariants = (managingVariantProduct.variants || []).map((v) =>
          v.id === variantId ? { ...v, ...updates } : v
        );
        const newTotalStock = nextVariants.reduce((sum, v) => sum + (v.stock || 0), 0);
        setManagingVariantProduct({
          ...managingVariantProduct,
          variants: nextVariants,
          total_stock: newTotalStock,
        });
      }

      showToast('Variant updated', 'success');
      await loadProducts();
    } catch (err: any) {
      console.error('Failed to update variant:', err);
      showToast(err.message || 'Failed to update variant', 'error');
    } finally {
      setVariantSubmitting(false);
    }
  };

  // Delete single variant from quick manager
  const handleDeleteQuickVariant = async (variantId: string) => {
    if (!managingVariantProduct || (managingVariantProduct.variants || []).length <= 1) {
      showToast('Cannot delete the only variant of a product', 'info');
      return;
    }
    try {
      await deleteProductVariant(variantId);
      const nextVariants = (managingVariantProduct.variants || []).filter((v) => v.id !== variantId);
      setManagingVariantProduct({
        ...managingVariantProduct,
        variants: nextVariants,
      });
      showToast('Variant deleted', 'success');
      await loadProducts();
    } catch (err: any) {
      showToast('Failed to delete variant', 'error');
    }
  };

  // Create new category
  const handleCreateCategory = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newCatName.trim()) return;
    setCatSubmitting(true);
    try {
      await createCategory({
        name: newCatName.trim(),
        description: newCatDescription.trim(),
        image_url: PRESET_IMAGES[0],
      });
      showToast(`Category "${newCatName}" added!`, 'success');
      setNewCatName('');
      setNewCatDescription('');
      setIsNewCatModalOpen(false);
      await loadCategories();
    } catch (err: any) {
      showToast(err.message || 'Failed to create category', 'error');
    } finally {
      setCatSubmitting(false);
    }
  };

  // Order status update
  const handleOrderStatusChange = async (orderId: string, newStatus: string) => {
    try {
      await updateOrderStatus(orderId, { status: newStatus });
      setOrders((prev) =>
        prev.map((o) => (o.id === orderId ? { ...o, status: newStatus } : o))
      );
      showToast(`Order status updated to ${newStatus}`, 'success');
    } catch (err) {
      showToast('Failed to update order status', 'error');
    }
  };

  const handleConfirmDeleteOrder = async () => {
    if (!deletingOrder) return;
    setIsDeleting(true);
    try {
      await deleteOrder(deletingOrder.id);
      setOrders((prev) => prev.filter((o) => o.id !== deletingOrder.id));
      showToast('Order deleted', 'success');
      setDeletingOrder(null);
    } catch (err) {
      showToast('Failed to delete order', 'error');
    } finally {
      setIsDeleting(false);
    }
  };

  // Filtered & Sorted Products computation
  const filteredProducts = useMemo(() => {
    let result = [...products];

    // Search query
    if (searchQuery.trim()) {
      const q = searchQuery.toLowerCase();
      result = result.filter(
        (p) =>
          p.name.toLowerCase().includes(q) ||
          (p.slug && p.slug.toLowerCase().includes(q)) ||
          p.description.toLowerCase().includes(q) ||
          p.variants?.some((v) => v.sku?.toLowerCase().includes(q) || v.color?.toLowerCase().includes(q))
      );
    }

    // Category filter
    if (selectedCategory !== 'all') {
      result = result.filter(
        (p) =>
          p.category_id === selectedCategory ||
          p.category.toLowerCase() === selectedCategory.toLowerCase()
      );
    }

    // Status filter
    if (selectedStatus !== 'all') {
      result = result.filter((p) => (p.status || 'active').toLowerCase() === selectedStatus.toLowerCase());
    }

    // Stock level filter
    if (stockFilter === 'out_of_stock') {
      result = result.filter((p) => (p.total_stock || 0) === 0);
    } else if (stockFilter === 'low_stock') {
      result = result.filter((p) => (p.total_stock || 0) > 0 && (p.total_stock || 0) <= 15);
    } else if (stockFilter === 'in_stock') {
      result = result.filter((p) => (p.total_stock || 0) > 15);
    }

    // Price range filter
    const min = minPrice !== '' ? Number(minPrice) : null;
    const max = maxPrice !== '' ? Number(maxPrice) : null;
    if (min !== null && !isNaN(min)) {
      result = result.filter((p) => (p.rawPrice || 0) >= min);
    }
    if (max !== null && !isNaN(max)) {
      result = result.filter((p) => (p.rawPrice || 0) <= max);
    }

    // Sorting
    result.sort((a, b) => {
      if (sortBy === 'name_asc') return a.name.localeCompare(b.name);
      if (sortBy === 'name_desc') return b.name.localeCompare(a.name);
      if (sortBy === 'price_asc') return (a.rawPrice || 0) - (b.rawPrice || 0);
      if (sortBy === 'price_desc') return (b.rawPrice || 0) - (a.rawPrice || 0);
      if (sortBy === 'stock_desc') return (b.total_stock || 0) - (a.total_stock || 0);
      return 0; // default newest/id order
    });

    return result;
  }, [products, searchQuery, selectedCategory, selectedStatus, stockFilter, sortBy, minPrice, maxPrice]);

  // Aggregate Metrics
  const totalStockUnits = useMemo(() => {
    return products.reduce((acc, p) => acc + (p.total_stock || 0), 0);
  }, [products]);

  const activeProductCount = useMemo(() => {
    return products.filter((p) => (p.status || 'active') === 'active').length;
  }, [products]);

  const lowStockCount = useMemo(() => {
    return products.filter((p) => (p.total_stock || 0) <= 15).length;
  }, [products]);

  // If auth is still checking
  if (authLoading) {
    return (
      <div className="min-h-screen bg-primary pt-36 pb-20 flex flex-col items-center justify-center">
        <div className="w-10 h-10 border-2 border-secondary/20 border-t-secondary rounded-full animate-spin mb-4" />
        <p className="text-xs uppercase tracking-editorial text-dark/60">Verifying administrative access...</p>
      </div>
    );
  }

  // Admin Access Gate
  if (!isAuthenticated || !isAdmin) {
    return (
      <div className="min-h-screen bg-primary pt-36 pb-24 flex items-center justify-center px-4">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="max-w-md w-full bg-white/80 backdrop-blur-md p-8 border border-dark/10 shadow-xl text-center"
        >
          <div className="mx-auto w-14 h-14 rounded-full bg-secondary/15 flex items-center justify-center text-secondary mb-5">
            <ShieldCheck size={28} />
          </div>
          <h1 className="font-heading text-2xl text-dark mb-2">Admin Portal</h1>
          <p className="text-xs text-dark/70 font-light leading-relaxed mb-6">
            Catalog CRUD operations, inventory management, and customer orders require administrative authorization.
          </p>
          <button
            onClick={() => openAuthModal('login')}
            className="w-full bg-dark hover:bg-dark-soft text-primary py-3 text-xs uppercase tracking-editorial transition-colors mb-4 flex items-center justify-center gap-2"
          >
            Sign In with Admin Account
          </button>
          <div className="bg-primary/50 border border-dark/5 p-3 rounded text-[11px] text-dark/60 text-left">
            <span className="font-medium text-dark block mb-0.5">Admin Account on Record:</span>
            <code>aidan.rafa@gmail.com</code>
          </div>
          <Link
            to="/shop"
            className="mt-6 inline-flex items-center gap-2 text-[11px] uppercase tracking-editorial text-dark/50 hover:text-secondary transition-colors"
          >
            <ArrowLeft size={13} /> Return to storefront
          </Link>
        </motion.div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-primary pt-28 lg:pt-36 pb-24 text-dark font-body">
      {/* Toast Notification */}
      <AnimatePresence>
        {toast && (
          <motion.div
            initial={{ opacity: 0, y: -20, x: 20 }}
            animate={{ opacity: 1, y: 0, x: 0 }}
            exit={{ opacity: 0, y: -20 }}
            className={`fixed top-6 right-6 z-[120] flex items-center gap-3 px-4 py-3 shadow-2xl border text-xs tracking-wide ${
              toast.type === 'error'
                ? 'bg-red-900 text-white border-red-700'
                : toast.type === 'info'
                ? 'bg-dark text-white border-dark/30'
                : 'bg-dark text-white border-secondary/40'
            }`}
          >
            {toast.type === 'error' ? (
              <AlertCircle size={16} className="text-red-400 shrink-0" />
            ) : (
              <CheckCircle2 size={16} className="text-emerald-400 shrink-0" />
            )}
            <span>{toast.message}</span>
            <button onClick={() => setToast(null)} className="ml-2 text-white/60 hover:text-white">
              <X size={14} />
            </button>
          </motion.div>
        )}
      </AnimatePresence>

      <div className="container-editorial max-w-7xl">
        {/* Breadcrumb & Navigation */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
          <Link
            to="/shop"
            className="inline-flex items-center gap-2 text-[11px] uppercase tracking-editorial text-dark/60 transition-colors hover:text-secondary"
          >
            <ArrowLeft size={14} /> Back to shop
          </Link>
          <div className="flex items-center gap-2 text-[11px] uppercase tracking-editorial text-dark/50">
            <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
            Backend Active: <span className="font-medium text-dark">localhost:8090</span>
          </div>
        </div>

        {/* Header Title Bar */}
        <div className="flex flex-col md:flex-row md:items-end justify-between gap-4 mb-8 pb-6 border-b border-dark/10">
          <div>
            <div className="flex items-center gap-2 text-[11px] uppercase tracking-editorial text-secondary font-medium mb-1">
              <ShieldCheck size={14} /> Katarina Daily Backoffice
            </div>
            <h1 className="font-heading text-3xl sm:text-4xl text-dark tracking-tight">Admin Dashboard</h1>
            <p className="text-xs text-dark/60 mt-1 font-light">
              Full control over catalog products, stock variants, imagery, and store orders.
            </p>
          </div>

          {/* Tab Navigation */}
          <div className="flex items-center gap-2 bg-white/60 p-1 border border-dark/10 rounded-sm">
            <button
              onClick={() => setActiveTab('products')}
              className={`flex items-center gap-2 px-4 py-2 text-xs uppercase tracking-editorial transition-all ${
                activeTab === 'products'
                  ? 'bg-dark text-white shadow-sm font-medium'
                  : 'text-dark/60 hover:text-dark hover:bg-dark/5'
              }`}
            >
              <Package size={14} />
              Products ({products.length})
            </button>
            <button
              onClick={() => setActiveTab('orders')}
              className={`flex items-center gap-2 px-4 py-2 text-xs uppercase tracking-editorial transition-all ${
                activeTab === 'orders'
                  ? 'bg-dark text-white shadow-sm font-medium'
                  : 'text-dark/60 hover:text-dark hover:bg-dark/5'
              }`}
            >
              <ShoppingBag size={14} />
              Orders ({orders.length})
            </button>
            <button
              onClick={() => setActiveTab('categories')}
              className={`flex items-center gap-2 px-4 py-2 text-xs uppercase tracking-editorial transition-all ${
                activeTab === 'categories'
                  ? 'bg-dark text-white shadow-sm font-medium'
                  : 'text-dark/60 hover:text-dark hover:bg-dark/5'
              }`}
            >
              <Layers size={14} />
              Categories ({categories.length})
            </button>
          </div>
        </div>

        {/* ========================================================= */}
        {/* TAB 1: PRODUCT MANAGEMENT (CRUD)                          */}
        {/* ========================================================= */}
        {activeTab === 'products' && (
          <div className="space-y-6">
            {/* KPI Metric Cards */}
            <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
              <div className="bg-white/70 backdrop-blur border border-dark/10 p-5 rounded-sm">
                <div className="flex items-center justify-between text-dark/50 text-[11px] uppercase tracking-editorial mb-2">
                  <span>Total Products</span>
                  <Package size={16} className="text-secondary" />
                </div>
                <div className="font-heading text-3xl text-dark">{products.length}</div>
                <div className="text-[10px] text-dark/50 mt-1">Catalog items listed</div>
              </div>

              <div className="bg-white/70 backdrop-blur border border-dark/10 p-5 rounded-sm">
                <div className="flex items-center justify-between text-dark/50 text-[11px] uppercase tracking-editorial mb-2">
                  <span>Active in Store</span>
                  <Sparkles size={16} className="text-emerald-600" />
                </div>
                <div className="font-heading text-3xl text-dark">{activeProductCount}</div>
                <div className="text-[10px] text-emerald-700 mt-1 font-medium">Publicly purchasable</div>
              </div>

              <div className="bg-white/70 backdrop-blur border border-dark/10 p-5 rounded-sm">
                <div className="flex items-center justify-between text-dark/50 text-[11px] uppercase tracking-editorial mb-2">
                  <span>Total Stock</span>
                  <Boxes size={16} className="text-secondary" />
                </div>
                <div className="font-heading text-3xl text-dark">{totalStockUnits}</div>
                <div className="text-[10px] text-dark/50 mt-1">Physical units across variants</div>
              </div>

              <div className="bg-white/70 backdrop-blur border border-dark/10 p-5 rounded-sm">
                <div className="flex items-center justify-between text-dark/50 text-[11px] uppercase tracking-editorial mb-2">
                  <span>Low / Out Stock</span>
                  <AlertTriangle size={16} className={lowStockCount > 0 ? 'text-amber-600' : 'text-dark/40'} />
                </div>
                <div className={`font-heading text-3xl ${lowStockCount > 0 ? 'text-amber-700' : 'text-dark'}`}>
                  {lowStockCount}
                </div>
                <div className="text-[10px] text-dark/50 mt-1">Stock ≤ 15 units</div>
              </div>
            </div>

            {/* Toolbar: Search, Filters & Action Button */}
            <div className="bg-white/80 backdrop-blur border border-dark/10 p-4 rounded-sm space-y-4">
              <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
                {/* Search Bar */}
                <div className="relative flex-1 max-w-md">
                  <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-dark/40" />
                  <input
                    type="text"
                    placeholder="Search by title, SKU, slug, color..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="w-full bg-primary/40 border border-dark/15 pl-9 pr-8 py-2 text-xs text-dark placeholder:text-dark/40 focus:outline-none focus:border-secondary transition-colors"
                  />
                  {searchQuery && (
                    <button
                      onClick={() => setSearchQuery('')}
                      className="absolute right-2.5 top-1/2 -translate-y-1/2 text-dark/40 hover:text-dark"
                    >
                      <X size={14} />
                    </button>
                  )}
                </div>

                {/* Primary Action Button: Add Product */}
                <div className="flex items-center gap-3">
                  <button
                    onClick={loadProducts}
                    title="Refresh product list"
                    className="p-2 border border-dark/15 hover:bg-dark/5 text-dark/70 hover:text-dark transition-colors"
                  >
                    <RefreshCw size={15} className={loadingProducts ? 'animate-spin' : ''} />
                  </button>

                  <button
                    onClick={handleOpenCreateModal}
                    className="flex items-center gap-2 bg-dark hover:bg-dark-soft text-primary px-5 py-2.5 text-xs uppercase tracking-editorial shadow-md transition-all active:scale-[0.98]"
                  >
                    <Plus size={15} /> Add New Product
                  </button>
                </div>
              </div>

              {/* Filter Controls Row */}
              <div className="flex flex-wrap items-center gap-3 pt-3 border-t border-dark/10 text-xs">
                {/* Category Filter */}
                <div className="flex items-center gap-1.5">
                  <span className="text-[10px] uppercase tracking-editorial text-dark/50">Category:</span>
                  <select
                    value={selectedCategory}
                    onChange={(e) => setSelectedCategory(e.target.value)}
                    className="bg-primary/30 border border-dark/15 px-2.5 py-1 text-xs outline-none focus:border-secondary"
                  >
                    <option value="all">All Categories</option>
                    {categories.map((c) => (
                      <option key={c.id} value={c.id}>
                        {c.name}
                      </option>
                    ))}
                  </select>
                </div>

                {/* Status Filter */}
                <div className="flex items-center gap-1.5">
                  <span className="text-[10px] uppercase tracking-editorial text-dark/50">Status:</span>
                  <select
                    value={selectedStatus}
                    onChange={(e) => setSelectedStatus(e.target.value)}
                    className="bg-primary/30 border border-dark/15 px-2.5 py-1 text-xs outline-none focus:border-secondary"
                  >
                    <option value="all">All Statuses</option>
                    <option value="active">Active (Visible)</option>
                    <option value="draft">Draft (Hidden)</option>
                    <option value="inactive">Inactive</option>
                  </select>
                </div>

                {/* Stock Filter */}
                <div className="flex items-center gap-1.5">
                  <span className="text-[10px] uppercase tracking-editorial text-dark/50">Stock:</span>
                  <select
                    value={stockFilter}
                    onChange={(e) => setStockFilter(e.target.value as any)}
                    className="bg-primary/30 border border-dark/15 px-2.5 py-1 text-xs outline-none focus:border-secondary"
                  >
                    <option value="all">All Inventory</option>
                    <option value="in_stock">In Stock (&gt; 15)</option>
                    <option value="low_stock">Low Stock (1-15)</option>
                    <option value="out_of_stock">Out of Stock (0)</option>
                  </select>
                </div>

                {/* Price Range Filter */}
                <div className="flex items-start gap-1.5">
                  <span className="text-[10px] uppercase tracking-editorial text-dark/50 mt-[7px]">Price:</span>
                  <div className="flex items-start gap-1">
                    {/* Min Price */}
                    <div className="flex flex-col gap-0.5">
                      <div className="relative">
                        <span className="absolute left-2 top-1/2 -translate-y-1/2 text-[10px] text-dark/40 pointer-events-none">Rp</span>
                        <input
                          type="text"
                          inputMode="numeric"
                          placeholder="Min"
                          value={minPrice}
                          onChange={(e) => {
                            const raw = e.target.value;
                            setMinPrice(raw);
                            if (raw === '') {
                              setMinPriceError(null);
                              return;
                            }
                            const num = Number(raw);
                            if (isNaN(num) || raw.trim() === '') {
                              setMinPriceError('harga tidak valid');
                            } else if (num < 0) {
                              setMinPriceError('harga tidak boleh dibawah 0');
                            } else if (num === 0) {
                              setMinPriceError('harga tidak boleh 0');
                            } else {
                              setMinPriceError(null);
                            }
                          }}
                          className={`bg-primary/30 border pl-7 pr-2 py-1 text-xs outline-none w-28 transition-colors ${
                            minPriceError
                              ? 'border-red-400 focus:border-red-500 text-red-700'
                              : 'border-dark/15 focus:border-secondary'
                          }`}
                        />
                      </div>
                      {minPriceError && (
                        <span className="text-[10px] text-red-500 leading-tight whitespace-nowrap">
                          {minPriceError}
                        </span>
                      )}
                    </div>

                    <span className="text-dark/30 text-[10px] mt-[7px]">–</span>

                    {/* Max Price */}
                    <div className="relative">
                      <span className="absolute left-2 top-1/2 -translate-y-1/2 text-[10px] text-dark/40 pointer-events-none">Rp</span>
                      <input
                        type="number"
                        placeholder="Max"
                        value={maxPrice}
                        onChange={(e) => setMaxPrice(e.target.value)}
                        min={0}
                        className="bg-primary/30 border border-dark/15 pl-7 pr-2 py-1 text-xs outline-none focus:border-secondary w-28 [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
                      />
                    </div>

                    {(minPrice !== '' || maxPrice !== '') && (
                      <button
                        onClick={() => { setMinPrice(''); setMaxPrice(''); setMinPriceError(null); }}
                        title="Clear price filter"
                        className="text-dark/40 hover:text-secondary transition-colors mt-[5px]"
                      >
                        <X size={13} />
                      </button>
                    )}
                  </div>
                </div>

                {/* Sort Filter */}
                <div className="flex items-center gap-1.5 ml-auto">
                  <span className="text-[10px] uppercase tracking-editorial text-dark/50">Sort:</span>
                  <select
                    value={sortBy}
                    onChange={(e) => setSortBy(e.target.value as any)}
                    className="bg-primary/30 border border-dark/15 px-2.5 py-1 text-xs outline-none focus:border-secondary"
                  >
                    <option value="newest">Default Order</option>
                    <option value="name_asc">Name: A to Z</option>
                    <option value="name_desc">Name: Z to A</option>
                    <option value="price_asc">Price: Low to High</option>
                    <option value="price_desc">Price: High to Low</option>
                    <option value="stock_desc">Highest Stock</option>
                  </select>
                </div>
              </div>
            </div>

            {/* Products Table */}
            <div className="bg-white/80 backdrop-blur border border-dark/10 shadow-sm overflow-hidden">
              {loadingProducts ? (
                <div className="p-12 text-center">
                  <div className="w-8 h-8 border-2 border-secondary/20 border-t-secondary rounded-full animate-spin mx-auto mb-3" />
                  <p className="text-xs uppercase tracking-editorial text-dark/50">Loading catalog products...</p>
                </div>
              ) : filteredProducts.length === 0 ? (
                <div className="p-16 text-center">
                  <Package size={36} className="text-dark/20 mx-auto mb-3" />
                  <h3 className="font-heading text-lg text-dark mb-1">No products found</h3>
                  <p className="text-xs text-dark/50 max-w-sm mx-auto mb-6">
                    {searchQuery || selectedCategory !== 'all' || selectedStatus !== 'all'
                      ? 'No items match your active filters. Try adjusting your query or resetting filters.'
                      : 'Your catalog has no products yet. Click "Add New Product" to create your first item.'}
                  </p>
                  <div className="flex justify-center gap-3">
                    {(searchQuery || selectedCategory !== 'all' || selectedStatus !== 'all') && (
                      <button
                        onClick={() => {
                          setSearchQuery('');
                          setSelectedCategory('all');
                          setSelectedStatus('all');
                          setStockFilter('all');
                        }}
                        className="px-4 py-2 border border-dark/20 text-xs uppercase tracking-editorial hover:bg-dark/5"
                      >
                        Reset Filters
                      </button>
                    )}
                    <button
                      onClick={handleOpenCreateModal}
                      className="px-4 py-2 bg-dark text-white text-xs uppercase tracking-editorial hover:bg-dark-soft"
                    >
                      + Add Product
                    </button>
                  </div>
                </div>
              ) : (
                <div className="overflow-x-auto">
                  <table className="w-full text-left border-collapse">
                    <thead>
                      <tr className="border-b border-dark/10 bg-primary/40 text-[10px] uppercase tracking-editorial text-dark/60">
                        <th className="py-3 px-4 font-medium">Product &amp; Category</th>
                        <th className="py-3 px-4 font-medium">Status</th>
                        <th className="py-3 px-4 font-medium">Price</th>
                        <th className="py-3 px-4 font-medium">Inventory &amp; Variants</th>
                        <th className="py-3 px-4 font-medium text-right">Actions</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-dark/5 text-xs">
                      {filteredProducts.map((p) => {
                        const variantCount = p.variants?.length || 0;
                        const stock = p.total_stock || 0;
                        const status = (p.status || 'active').toLowerCase();

                        return (
                          <tr key={p.id} className="hover:bg-primary/20 transition-colors group">
                            {/* Product Info */}
                            <td className="py-4 px-4">
                              <div className="flex items-center gap-3.5">
                                <div className="w-12 h-14 bg-dark/5 shrink-0 border border-dark/10 overflow-hidden relative">
                                  <img
                                    src={p.image}
                                    alt={p.name}
                                    className="w-full h-full object-cover object-center group-hover:scale-105 transition-transform duration-300"
                                    onError={(e) => {
                                      (e.target as HTMLImageElement).src = PRESET_IMAGES[0];
                                    }}
                                  />
                                </div>
                                <div className="min-w-0">
                                  <div className="flex items-center gap-2">
                                    <h4 className="font-medium text-dark text-xs sm:text-sm truncate max-w-xs sm:max-w-md">
                                      {p.name}
                                    </h4>
                                    <span className="text-[10px] px-2 py-0.5 bg-primary/80 border border-dark/10 rounded-full text-dark/70 shrink-0">
                                      {p.category}
                                    </span>
                                  </div>
                                  <div className="text-[11px] text-dark/40 font-mono mt-0.5 truncate max-w-xs">
                                    /{p.slug || p.id}
                                  </div>
                                </div>
                              </div>
                            </td>

                            {/* Status */}
                            <td className="py-4 px-4">
                              <span
                                className={`inline-flex items-center gap-1.5 px-2.5 py-1 text-[10px] uppercase tracking-editorial rounded-full font-medium ${
                                  status === 'active'
                                    ? 'bg-emerald-50 text-emerald-700 border border-emerald-200'
                                    : status === 'draft'
                                    ? 'bg-amber-50 text-amber-700 border border-amber-200'
                                    : 'bg-rose-50 text-rose-700 border border-rose-200'
                                }`}
                              >
                                <span
                                  className={`w-1.5 h-1.5 rounded-full ${
                                    status === 'active'
                                      ? 'bg-emerald-500'
                                      : status === 'draft'
                                      ? 'bg-amber-500'
                                      : 'bg-rose-500'
                                  }`}
                                />
                                {status}
                              </span>
                            </td>

                            {/* Price */}
                            <td className="py-4 px-4">
                              <span className="font-medium text-dark">{p.price}</span>
                            </td>

                            {/* Inventory & Variants */}
                            <td className="py-4 px-4">
                              <div className="space-y-1">
                                <div className="flex items-center gap-2">
                                  <span
                                    className={`font-semibold ${
                                      stock === 0
                                        ? 'text-red-600'
                                        : stock <= 15
                                        ? 'text-amber-600'
                                        : 'text-dark'
                                    }`}
                                  >
                                    {stock} in stock
                                  </span>
                                  <span className="text-dark/30 text-[10px]">•</span>
                                  <button
                                    onClick={() => setManagingVariantProduct(p)}
                                    className="text-[11px] text-secondary hover:underline flex items-center gap-1"
                                    title="Click to view/manage variants"
                                  >
                                    <Layers size={11} /> {variantCount} variants
                                  </button>
                                </div>

                                {/* Mini color/size pills preview */}
                                {p.variants && p.variants.length > 0 && (
                                  <div className="flex flex-wrap gap-1 max-w-xs">
                                    {p.variants.slice(0, 4).map((v) => (
                                      <span
                                        key={v.id || v.sku}
                                        className="text-[9px] bg-primary/60 border border-dark/10 px-1.5 py-0.2 rounded text-dark/70"
                                      >
                                        {v.color} ({v.size}): {v.stock}
                                      </span>
                                    ))}
                                    {p.variants.length > 4 && (
                                      <span className="text-[9px] text-dark/40 self-center">
                                        +{p.variants.length - 4} more
                                      </span>
                                    )}
                                  </div>
                                )}
                              </div>
                            </td>

                            {/* Action Buttons */}
                            <td className="py-4 px-4 text-right">
                              <div className="inline-flex items-center gap-2">
                                <button
                                  onClick={() => setManagingVariantProduct(p)}
                                  className="p-1.5 text-dark/60 hover:text-secondary hover:bg-secondary/10 rounded transition-colors"
                                  title="Manage Variants &amp; Stock"
                                >
                                  <SlidersHorizontal size={15} />
                                </button>
                                <button
                                  onClick={() => handleOpenEditModal(p)}
                                  className="p-1.5 text-dark/60 hover:text-dark hover:bg-dark/5 rounded transition-colors"
                                  title="Edit Product Details"
                                >
                                  <Edit2 size={15} />
                                </button>
                                <Link
                                  to={`/product/${p.slug || p.id}`}
                                  target="_blank"
                                  className="p-1.5 text-dark/60 hover:text-dark hover:bg-dark/5 rounded transition-colors"
                                  title="View on storefront"
                                >
                                  <ExternalLink size={15} />
                                </Link>
                                <button
                                  onClick={() => setDeletingProduct(p)}
                                  className="p-1.5 text-red-500 hover:text-red-700 hover:bg-red-50 rounded transition-colors"
                                  title="Delete Product"
                                >
                                  <Trash2 size={15} />
                                </button>
                              </div>
                            </td>
                          </tr>
                        );
                      })}
                    </tbody>
                  </table>
                </div>
              )}

              {/* Footer Table Summary */}
              {!loadingProducts && filteredProducts.length > 0 && (
                <div className="px-4 py-3 bg-primary/30 border-t border-dark/10 flex items-center justify-between text-[11px] text-dark/60">
                  <span>
                    Showing <strong>{filteredProducts.length}</strong> of <strong>{products.length}</strong> catalog items
                  </span>
                  <span>Click variant links for quick inventory adjustments</span>
                </div>
              )}
            </div>
          </div>
        )}

        {/* ========================================================= */}
        {/* TAB 2: ORDERS MANAGEMENT                                  */}
        {/* ========================================================= */}
        {activeTab === 'orders' && (
          <div className="space-y-6">
            <div className="flex justify-between items-center bg-white/80 p-4 border border-dark/10">
              <div>
                <h2 className="font-heading text-xl text-dark">Order Management</h2>
                <p className="text-xs text-dark/50">Manage customer purchases and update order statuses.</p>
              </div>
              <button
                onClick={loadOrders}
                className="flex items-center gap-2 border border-dark/15 px-3 py-1.5 text-xs text-dark hover:bg-dark/5"
              >
                <RefreshCw size={13} className={loadingOrders ? 'animate-spin' : ''} /> Refresh
              </button>
            </div>

            {loadingOrders ? (
              <div className="p-12 text-center bg-white/70 border border-dark/10">
                <div className="w-8 h-8 border-2 border-secondary/20 border-t-secondary rounded-full animate-spin mx-auto mb-3" />
                <p className="text-xs uppercase tracking-editorial text-dark/50">Loading orders...</p>
              </div>
            ) : orders.length === 0 ? (
              <div className="p-12 text-center bg-white/70 border border-dark/10">
                <ShoppingBag size={32} className="text-dark/20 mx-auto mb-2" />
                <p className="text-xs text-dark/50">No orders found in database.</p>
              </div>
            ) : (
              <div className="overflow-x-auto bg-white/80 border border-dark/10 shadow-sm">
                <table className="w-full text-left border-collapse">
                  <thead>
                    <tr className="border-b border-dark/10 bg-primary/40 text-[10px] uppercase tracking-editorial text-dark/60">
                      <th className="px-4 py-3">Order Number</th>
                      <th className="px-4 py-3">Recipient</th>
                      <th className="px-4 py-3">Total Amount</th>
                      <th className="px-4 py-3">Status</th>
                      <th className="px-4 py-3 text-right">Actions</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-dark/5 text-xs">
                    {orders.map((order) => (
                      <tr key={order.id} className="hover:bg-primary/20 transition-colors">
                        <td className="px-4 py-3.5 font-mono text-dark">{order.order_number}</td>
                        <td className="px-4 py-3.5">
                          <div className="font-medium text-dark">{order.shipping_recipient_name}</div>
                          <div className="text-[10px] text-dark/50">{order.shipping_city || order.shipping_phone}</div>
                        </td>
                        <td className="px-4 py-3.5 font-medium">{formatIDR(order.total)}</td>
                        <td className="px-4 py-3.5">
                          <select
                            value={order.status}
                            onChange={(e) => handleOrderStatusChange(order.id, e.target.value)}
                            className="rounded border border-dark/20 bg-white px-2.5 py-1 text-xs outline-none focus:border-secondary font-medium"
                          >
                            {['paid', 'processing', 'shipped', 'completed', 'cancelled'].map((s) => (
                              <option key={s} value={s}>
                                {s.toUpperCase()}
                              </option>
                            ))}
                          </select>
                        </td>
                        <td className="px-4 py-3.5 text-right space-x-2">
                          <button
                            onClick={() => setDeletingOrder(order)}
                            className="text-red-500 hover:text-red-700 text-xs px-2 py-1"
                          >
                            Delete
                          </button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}

        {/* ========================================================= */}
        {/* TAB 3: CATEGORIES MANAGEMENT                              */}
        {/* ========================================================= */}
        {activeTab === 'categories' && (
          <div className="space-y-6">
            <div className="flex justify-between items-center bg-white/80 p-4 border border-dark/10">
              <div>
                <h2 className="font-heading text-xl text-dark">Category Management</h2>
                <p className="text-xs text-dark/50">Manage apparel classifications and taxonomies.</p>
              </div>
              <button
                onClick={() => setIsNewCatModalOpen(true)}
                className="flex items-center gap-2 bg-dark text-white px-4 py-2 text-xs uppercase tracking-editorial"
              >
                <Plus size={14} /> Add Category
              </button>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
              {categories.map((c) => {
                const count = products.filter(
                  (p) => p.category_id === c.id || p.category.toLowerCase() === c.name.toLowerCase()
                ).length;
                return (
                  <div key={c.id} className="bg-white/80 border border-dark/10 p-5 rounded-sm space-y-3">
                    <div className="flex items-center justify-between">
                      <h3 className="font-heading text-lg text-dark">{c.name}</h3>
                      <span className="text-[10px] bg-primary px-2 py-0.5 rounded text-dark/70 font-mono">
                        /{c.slug}
                      </span>
                    </div>
                    <p className="text-xs text-dark/60 font-light line-clamp-2">
                      {c.description || 'No description provided.'}
                    </p>
                    <div className="pt-2 border-t border-dark/10 flex items-center justify-between text-xs text-dark/50">
                      <span>Products: <strong>{count}</strong></span>
                      <span className="font-mono text-[9px] text-dark/30 truncate max-w-[100px]">{c.id}</span>
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        )}
      </div>

      {/* ========================================================= */}
      {/* MODAL 1: CREATE / EDIT PRODUCT FORM (FULL CRUD)           */}
      {/* ========================================================= */}
      <AnimatePresence>
        {isFormModalOpen && (
          <div className="fixed inset-0 z-[110] flex items-center justify-center p-4 sm:p-6 overflow-y-auto">
            {/* Backdrop */}
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => !formSubmitting && setIsFormModalOpen(false)}
              className="fixed inset-0 bg-dark/60 backdrop-blur-sm"
            />

            {/* Modal Body */}
            <motion.div
              initial={{ opacity: 0, scale: 0.95, y: 15 }}
              animate={{ opacity: 1, scale: 1, y: 0 }}
              exit={{ opacity: 0, scale: 0.95, y: 15 }}
              className="relative w-full max-w-3xl bg-primary-100 border border-dark/20 shadow-2xl z-10 max-h-[90vh] flex flex-col"
            >
              {/* Modal Header */}
              <div className="flex items-center justify-between p-6 border-b border-dark/10 bg-white/70">
                <div>
                  <h2 className="font-heading text-2xl text-dark">
                    {editingProduct ? 'Edit Product' : 'Create New Product'}
                  </h2>
                  <p className="text-xs text-dark/60 mt-0.5">
                    {editingProduct
                      ? 'Update product details, pricing, media, and inventory variants.'
                      : 'Add a new piece to the Katarina Daily catalog with color & size variants.'}
                  </p>
                </div>
                <button
                  disabled={formSubmitting}
                  onClick={() => setIsFormModalOpen(false)}
                  className="p-1.5 text-dark/40 hover:text-dark rounded transition-colors"
                >
                  <X size={20} />
                </button>
              </div>

              {/* Form Content */}
              <form onSubmit={handleSaveProduct} className="flex-1 overflow-y-auto p-6 space-y-6">
                {formError && (
                  <div className="p-3 bg-red-100 border border-red-300 text-red-800 text-xs flex items-center gap-2">
                    <AlertCircle size={16} className="shrink-0 text-red-600" />
                    <span>{formError}</span>
                  </div>
                )}

                {/* Section A: Basic Metadata */}
                <div className="bg-white/80 p-5 border border-dark/10 rounded-sm space-y-4">
                  <h3 className="text-xs font-semibold uppercase tracking-editorial text-secondary flex items-center gap-2">
                    <Tag size={14} /> 1. Basic Information
                  </h3>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    {/* Name */}
                    <div className="md:col-span-2">
                      <label className="block text-[10px] uppercase tracking-editorial text-dark/70 mb-1">
                        Product Title *
                      </label>
                      <input
                        type="text"
                        required
                        placeholder="e.g. Minimalist Linen Blouse"
                        value={formData.name}
                        onChange={(e) => handleNameChange(e.target.value)}
                        className="w-full bg-white border border-dark/20 p-2.5 text-xs text-dark outline-none focus:border-secondary"
                      />
                    </div>

                    {/* Slug */}
                    <div>
                      <label className="block text-[10px] uppercase tracking-editorial text-dark/70 mb-1">
                        URL Slug *
                      </label>
                      <input
                        type="text"
                        required
                        placeholder="e.g. minimalist-linen-blouse"
                        value={formData.slug}
                        onChange={(e) => setFormData({ ...formData, slug: e.target.value })}
                        className="w-full bg-white border border-dark/20 p-2.5 text-xs font-mono text-dark outline-none focus:border-secondary"
                      />
                    </div>

                    {/* Category */}
                    <div>
                      <label className="block text-[10px] uppercase tracking-editorial text-dark/70 mb-1">
                        Category *
                      </label>
                      <select
                        required
                        value={formData.category_id}
                        onChange={(e) => setFormData({ ...formData, category_id: e.target.value })}
                        className="w-full bg-white border border-dark/20 p-2.5 text-xs text-dark outline-none focus:border-secondary"
                      >
                        {categories.map((c) => (
                          <option key={c.id} value={c.id}>
                            {c.name}
                          </option>
                        ))}
                      </select>
                    </div>

                    {/* Status */}
                    <div>
                      <label className="block text-[10px] uppercase tracking-editorial text-dark/70 mb-1">
                        Visibility Status *
                      </label>
                      <select
                        value={formData.status}
                        onChange={(e) => setFormData({ ...formData, status: e.target.value as any })}
                        className="w-full bg-white border border-dark/20 p-2.5 text-xs text-dark outline-none focus:border-secondary"
                      >
                        <option value="active">Active (Visible in store)</option>
                        <option value="draft">Draft (Hidden in store)</option>
                        <option value="inactive">Inactive</option>
                      </select>
                    </div>

                    {/* Description */}
                    <div className="md:col-span-2">
                      <label className="block text-[10px] uppercase tracking-editorial text-dark/70 mb-1">
                        Product Description
                      </label>
                      <textarea
                        rows={3}
                        placeholder="Detailed fabric description, fit details, inner embroidery notes..."
                        value={formData.description}
                        onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                        className="w-full bg-white border border-dark/20 p-2.5 text-xs text-dark outline-none focus:border-secondary leading-relaxed"
                      />
                    </div>
                  </div>
                </div>

                {/* Section B: Media & Imagery */}
                <div className="bg-white/80 p-5 border border-dark/10 rounded-sm space-y-4">
                  <h3 className="text-xs font-semibold uppercase tracking-editorial text-secondary flex items-center gap-2">
                    <Sparkles size={14} /> 2. Product Imagery
                  </h3>

                  <div className="flex flex-col sm:flex-row gap-4 items-start">
                    <div className="w-24 h-32 bg-dark/5 border border-dark/20 shrink-0 overflow-hidden relative rounded-sm">
                      <img
                        src={formData.primaryImageUrl}
                        alt="Preview"
                        className="w-full h-full object-cover object-center"
                        onError={(e) => {
                          (e.target as HTMLImageElement).src = PRESET_IMAGES[0];
                        }}
                      />
                      <span className="absolute bottom-0 inset-x-0 bg-dark/70 text-white text-[9px] text-center py-0.5 font-medium">
                        Preview
                      </span>
                    </div>

                    <div className="flex-1 space-y-2 w-full">
                      <label className="block text-[10px] uppercase tracking-editorial text-dark/70">
                        Primary Image URL *
                      </label>
                      <input
                        type="url"
                        required
                        placeholder="https://images.pexels.com/..."
                        value={formData.primaryImageUrl}
                        onChange={(e) => setFormData({ ...formData, primaryImageUrl: e.target.value })}
                        className="w-full bg-white border border-dark/20 p-2 text-xs text-dark outline-none focus:border-secondary"
                      />

                      <div className="pt-2">
                        <span className="text-[10px] uppercase tracking-editorial text-dark/40 block mb-1.5">
                          Or select from studio samples:
                        </span>
                        <div className="flex flex-wrap gap-2">
                          {PRESET_IMAGES.map((url, idx) => (
                            <button
                              key={idx}
                              type="button"
                              onClick={() => setFormData({ ...formData, primaryImageUrl: url })}
                              className={`w-10 h-10 border rounded overflow-hidden transition-all ${
                                formData.primaryImageUrl === url ? 'ring-2 ring-secondary border-transparent scale-105' : 'border-dark/20 opacity-70 hover:opacity-100'
                              }`}
                            >
                              <img src={url} alt={`Preset ${idx}`} className="w-full h-full object-cover" />
                            </button>
                          ))}
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                {/* Section C: Variants Management (Sizes, Colors, Stock, Price) */}
                <div className="bg-white/80 p-5 border border-dark/10 rounded-sm space-y-4">
                  <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
                    <h3 className="text-xs font-semibold uppercase tracking-editorial text-secondary flex items-center gap-2">
                      <Boxes size={14} /> 3. Inventory Variants (Colors &amp; Sizes)
                    </h3>
                    <div className="flex items-center gap-2">
                      <button
                        type="button"
                        onClick={() => handleGenerateSizes('Cream')}
                        className="text-[10px] text-dark/70 hover:text-dark border border-dark/20 px-2 py-1 hover:bg-dark/5"
                      >
                        + Generate S, M, L
                      </button>
                      <button
                        type="button"
                        onClick={handleAddVariantRow}
                        className="text-[10px] bg-secondary text-white px-2.5 py-1 hover:bg-secondary-600 flex items-center gap-1 font-medium"
                      >
                        <Plus size={11} /> Add Variant
                      </button>
                    </div>
                  </div>

                  <p className="text-[11px] text-dark/50">
                    Define colorways, sizes, specific pricing, and available units for each sku.
                  </p>

                  <div className="overflow-x-auto border border-dark/10">
                    <table className="w-full text-left text-xs bg-white">
                      <thead>
                        <tr className="bg-primary/40 border-b border-dark/10 text-[9px] uppercase tracking-editorial text-dark/60">
                          <th className="p-2">SKU</th>
                          <th className="p-2">Color</th>
                          <th className="p-2">Size</th>
                          <th className="p-2">Price (Rp)</th>
                          <th className="p-2">Stock</th>
                          <th className="p-2">Weight (g)</th>
                          <th className="p-2 text-center">Action</th>
                        </tr>
                      </thead>
                      <tbody className="divide-y divide-dark/5">
                        {formData.variants.map((v, i) => (
                          <tr key={i} className="hover:bg-primary/20">
                            <td className="p-2">
                              <input
                                type="text"
                                required
                                value={v.sku}
                                placeholder="SKU code"
                                onChange={(e) => handleUpdateVariantRow(i, 'sku', e.target.value)}
                                className="w-24 font-mono text-[11px] border border-dark/20 p-1 bg-white"
                              />
                            </td>
                            <td className="p-2">
                              <input
                                type="text"
                                required
                                value={v.color}
                                placeholder="Color"
                                onChange={(e) => handleUpdateVariantRow(i, 'color', e.target.value)}
                                className="w-20 text-[11px] border border-dark/20 p-1 bg-white"
                              />
                            </td>
                            <td className="p-2">
                              <select
                                value={v.size}
                                onChange={(e) => handleUpdateVariantRow(i, 'size', e.target.value)}
                                className="text-[11px] border border-dark/20 p-1 bg-white"
                              >
                                {['XS', 'S', 'M', 'L', 'XL', 'All Size'].map((sz) => (
                                  <option key={sz} value={sz}>
                                    {sz}
                                  </option>
                                ))}
                              </select>
                            </td>
                            <td className="p-2">
                              <input
                                type="number"
                                required
                                min="1000"
                                value={v.price}
                                onChange={(e) => handleUpdateVariantRow(i, 'price', Number(e.target.value))}
                                className="w-24 text-[11px] border border-dark/20 p-1 bg-white"
                              />
                            </td>
                            <td className="p-2">
                              <input
                                type="number"
                                required
                                min="0"
                                value={v.stock}
                                onChange={(e) => handleUpdateVariantRow(i, 'stock', Number(e.target.value))}
                                className="w-16 text-[11px] border border-dark/20 p-1 bg-white"
                              />
                            </td>
                            <td className="p-2">
                              <input
                                type="number"
                                required
                                min="10"
                                value={v.weight}
                                onChange={(e) => handleUpdateVariantRow(i, 'weight', Number(e.target.value))}
                                className="w-16 text-[11px] border border-dark/20 p-1 bg-white"
                              />
                            </td>
                            <td className="p-2 text-center">
                              <button
                                type="button"
                                onClick={() => handleRemoveVariantRow(i)}
                                className="text-red-500 hover:text-red-700 p-1"
                                title="Remove variant row"
                              >
                                <Trash2 size={13} />
                              </button>
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>

                {/* Modal Footer Controls */}
                <div className="pt-4 border-t border-dark/10 flex items-center justify-end gap-3 sticky bottom-0 bg-primary-100 py-2">
                  <button
                    type="button"
                    disabled={formSubmitting}
                    onClick={() => setIsFormModalOpen(false)}
                    className="px-5 py-2.5 border border-dark/20 text-dark text-xs uppercase tracking-editorial hover:bg-dark/5"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    disabled={formSubmitting}
                    className="px-7 py-2.5 bg-dark hover:bg-dark-soft text-primary text-xs uppercase tracking-editorial font-medium shadow-lg transition-all flex items-center gap-2"
                  >
                    {formSubmitting ? (
                      <>
                        <div className="w-3.5 h-3.5 border-2 border-white/20 border-t-white rounded-full animate-spin" />
                        Saving Product...
                      </>
                    ) : (
                      <>
                        <Check size={14} />
                        {editingProduct ? 'Save Changes' : 'Publish Product'}
                      </>
                    )}
                  </button>
                </div>
              </form>
            </motion.div>
          </div>
        )}
      </AnimatePresence>

      {/* ========================================================= */}
      {/* MODAL 2: QUICK VARIANT & STOCK MANAGER                    */}
      {/* ========================================================= */}
      <AnimatePresence>
        {managingVariantProduct && (
          <div className="fixed inset-0 z-[110] flex items-center justify-center p-4 overflow-y-auto">
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setManagingVariantProduct(null)}
              className="fixed inset-0 bg-dark/60 backdrop-blur-sm"
            />
            <motion.div
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              exit={{ opacity: 0, scale: 0.95 }}
              className="relative w-full max-w-2xl bg-white border border-dark/20 shadow-2xl p-6 z-10"
            >
              <div className="flex items-center justify-between pb-4 mb-4 border-b border-dark/10">
                <div className="flex items-center gap-3">
                  <img
                    src={managingVariantProduct.image}
                    alt={managingVariantProduct.name}
                    className="w-10 h-12 object-cover border border-dark/10"
                  />
                  <div>
                    <h3 className="font-heading text-xl text-dark">{managingVariantProduct.name}</h3>
                    <p className="text-xs text-dark/50">
                      Total Stock: <strong className="text-dark">{managingVariantProduct.total_stock}</strong> units
                    </p>
                  </div>
                </div>
                <button
                  onClick={() => setManagingVariantProduct(null)}
                  className="p-1.5 text-dark/40 hover:text-dark"
                >
                  <X size={18} />
                </button>
              </div>

              <div className="space-y-4">
                <p className="text-xs text-dark/60 font-light">
                  Directly adjust inventory units and pricing for each colorway/size.
                </p>

                <div className="divide-y divide-dark/10 border border-dark/10 rounded-sm max-h-72 overflow-y-auto">
                  {(managingVariantProduct.variants || []).map((v) => (
                    <div key={v.id} className="p-3 bg-primary/20 flex flex-wrap items-center justify-between gap-3 text-xs">
                      <div>
                        <div className="font-semibold text-dark">
                          {v.color} — {v.size}
                        </div>
                        <div className="text-[10px] text-dark/40 font-mono">{v.sku}</div>
                      </div>

                      <div className="flex items-center gap-4">
                        {/* Price adjust */}
                        <div className="flex items-center gap-1">
                          <span className="text-[10px] text-dark/50">Rp</span>
                          <input
                            type="number"
                            defaultValue={v.price}
                            onBlur={(e) => handleUpdateQuickVariant(v.id, { price: Number(e.target.value) })}
                            className="w-24 p-1 border border-dark/20 bg-white text-xs font-medium"
                          />
                        </div>

                        {/* Stock adjust */}
                        <div className="flex items-center gap-1">
                          <span className="text-[10px] text-dark/50">Stock:</span>
                          <input
                            type="number"
                            min="0"
                            defaultValue={v.stock}
                            onBlur={(e) => handleUpdateQuickVariant(v.id, { stock: Number(e.target.value) })}
                            className="w-16 p-1 border border-dark/20 bg-white text-xs font-semibold text-center"
                          />
                        </div>

                        <button
                          onClick={() => handleDeleteQuickVariant(v.id)}
                          className="text-red-500 hover:text-red-700 p-1"
                          title="Delete variant"
                        >
                          <Trash2 size={13} />
                        </button>
                      </div>
                    </div>
                  ))}
                </div>

                <div className="flex justify-between items-center pt-2">
                  <button
                    onClick={() => {
                      const p = managingVariantProduct;
                      setManagingVariantProduct(null);
                      handleOpenEditModal(p);
                    }}
                    className="text-xs text-secondary hover:underline flex items-center gap-1"
                  >
                    <Edit2 size={12} /> Open Full Product Editor
                  </button>
                  <button
                    onClick={() => setManagingVariantProduct(null)}
                    className="px-5 py-2 bg-dark text-white text-xs uppercase tracking-editorial"
                  >
                    Done
                  </button>
                </div>
              </div>
            </motion.div>
          </div>
        )}
      </AnimatePresence>

      {/* ========================================================= */}
      {/* MODAL 3: DELETE CONFIRMATION                              */}
      {/* ========================================================= */}
      <AnimatePresence>
        {deletingProduct && (
          <div className="fixed inset-0 z-[120] flex items-center justify-center p-4">
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => !isDeleting && setDeletingProduct(null)}
              className="fixed inset-0 bg-dark/60 backdrop-blur-sm"
            />
            <motion.div
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              exit={{ opacity: 0, scale: 0.95 }}
              className="relative w-full max-w-md bg-white border border-dark/20 p-6 z-10 text-center space-y-4"
            >
              <div className="w-12 h-12 bg-red-100 text-red-600 rounded-full flex items-center justify-center mx-auto">
                <AlertTriangle size={24} />
              </div>
              <h3 className="font-heading text-xl text-dark">Delete Product?</h3>
              <p className="text-xs text-dark/70 font-light leading-relaxed">
                Are you sure you want to permanently delete <strong>"{deletingProduct.name}"</strong>? All associated
                variants, images, and catalog references will be deleted.
              </p>
              <div className="flex justify-center gap-3 pt-2">
                <button
                  disabled={isDeleting}
                  onClick={() => setDeletingProduct(null)}
                  className="px-5 py-2.5 border border-dark/20 text-dark text-xs uppercase tracking-editorial hover:bg-dark/5"
                >
                  Cancel
                </button>
                <button
                  disabled={isDeleting}
                  onClick={handleConfirmDeleteProduct}
                  className="px-5 py-2.5 bg-red-600 hover:bg-red-700 text-white text-xs uppercase tracking-editorial flex items-center gap-2"
                >
                  {isDeleting ? 'Deleting...' : 'Delete Permanently'}
                </button>
              </div>
            </motion.div>
          </div>
        )}
      </AnimatePresence>

      {/* Delete Order Confirmation */}
      <AnimatePresence>
        {deletingOrder && (
          <div className="fixed inset-0 z-[120] flex items-center justify-center p-4">
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => !isDeleting && setDeletingOrder(null)}
              className="fixed inset-0 bg-dark/60 backdrop-blur-sm"
            />
            <motion.div
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              exit={{ opacity: 0, scale: 0.95 }}
              className="relative w-full max-w-md bg-white border border-dark/20 p-6 z-10 text-center space-y-4"
            >
              <div className="w-12 h-12 bg-red-100 text-red-600 rounded-full flex items-center justify-center mx-auto">
                <AlertTriangle size={24} />
              </div>
              <h3 className="font-heading text-xl text-dark">Delete Order?</h3>
              <p className="text-xs text-dark/70">
                Are you sure you want to delete order <strong>#{deletingOrder.order_number}</strong>?
              </p>
              <div className="flex justify-center gap-3 pt-2">
                <button
                  onClick={() => setDeletingOrder(null)}
                  className="px-5 py-2 border border-dark/20 text-dark text-xs uppercase tracking-editorial"
                >
                  Cancel
                </button>
                <button
                  onClick={handleConfirmDeleteOrder}
                  className="px-5 py-2 bg-red-600 text-white text-xs uppercase tracking-editorial"
                >
                  {isDeleting ? 'Deleting...' : 'Delete Order'}
                </button>
              </div>
            </motion.div>
          </div>
        )}
      </AnimatePresence>

      {/* ========================================================= */}
      {/* MODAL 4: NEW CATEGORY MODAL                               */}
      {/* ========================================================= */}
      <AnimatePresence>
        {isNewCatModalOpen && (
          <div className="fixed inset-0 z-[120] flex items-center justify-center p-4">
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => !catSubmitting && setIsNewCatModalOpen(false)}
              className="fixed inset-0 bg-dark/60 backdrop-blur-sm"
            />
            <motion.div
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              exit={{ opacity: 0, scale: 0.95 }}
              className="relative w-full max-w-md bg-white border border-dark/20 p-6 z-10 space-y-4"
            >
              <div className="flex justify-between items-center border-b border-dark/10 pb-3">
                <h3 className="font-heading text-xl text-dark">Add New Category</h3>
                <button onClick={() => setIsNewCatModalOpen(false)} className="text-dark/40 hover:text-dark">
                  <X size={16} />
                </button>
              </div>

              <form onSubmit={handleCreateCategory} className="space-y-4 text-xs">
                <div>
                  <label className="block text-[10px] uppercase tracking-editorial text-dark/70 mb-1">
                    Category Name *
                  </label>
                  <input
                    type="text"
                    required
                    placeholder="e.g. Knitwear, Accessories"
                    value={newCatName}
                    onChange={(e) => setNewCatName(e.target.value)}
                    className="w-full border border-dark/20 p-2 text-xs outline-none focus:border-secondary"
                  />
                </div>

                <div>
                  <label className="block text-[10px] uppercase tracking-editorial text-dark/70 mb-1">
                    Description
                  </label>
                  <textarea
                    rows={2}
                    placeholder="Category details..."
                    value={newCatDescription}
                    onChange={(e) => setNewCatDescription(e.target.value)}
                    className="w-full border border-dark/20 p-2 text-xs outline-none focus:border-secondary"
                  />
                </div>

                <div className="pt-2 flex justify-end gap-3">
                  <button
                    type="button"
                    onClick={() => setIsNewCatModalOpen(false)}
                    className="px-4 py-2 border border-dark/20 text-xs uppercase tracking-editorial"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    disabled={catSubmitting}
                    className="px-5 py-2 bg-dark text-white text-xs uppercase tracking-editorial"
                  >
                    {catSubmitting ? 'Saving...' : 'Create Category'}
                  </button>
                </div>
              </form>
            </motion.div>
          </div>
        )}
      </AnimatePresence>
    </div>
  );
}
