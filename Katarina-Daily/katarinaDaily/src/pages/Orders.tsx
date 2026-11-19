import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { Package, ArrowLeft, Calendar, Truck, Clock } from 'lucide-react';
import { getOrders, type Order } from '@/services/orderService';
import { useAuth } from '@/context/AuthContext';
import { formatIDR } from '@/services/catalogService';

export default function Orders() {
  const { isAuthenticated, openAuthModal } = useAuth();
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (isAuthenticated) {
      setLoading(true);
      getOrders().then((list) => {
        setOrders(list);
        setLoading(false);
      });
    } else {
      setLoading(false);
    }
  }, [isAuthenticated]);

  if (!isAuthenticated) {
    return (
      <div className="bg-primary pt-32 lg:pt-40 pb-24 min-h-screen">
        <div className="container-editorial max-w-md text-center">
          <Package size={44} className="mx-auto text-dark/30 mb-4 stroke-1" />
          <h2 className="font-heading text-2xl text-dark">Sign in to view orders</h2>
          <p className="mt-2 text-xs font-light text-dark/60">
            Please log in to your Katarina Daily account to view your past and ongoing orders.
          </p>
          <button
            onClick={() => openAuthModal('login')}
            className="mt-6 bg-secondary px-8 py-3.5 text-xs uppercase tracking-editorial text-primary hover:bg-secondary-600 transition-colors"
          >
            Sign In
          </button>
        </div>
      </div>
    );
  }

  const getStatusBadge = (status: string) => {
    switch (status.toLowerCase()) {
      case 'paid':
      case 'processing':
        return 'bg-blue-50 text-blue-700 border-blue-200';
      case 'shipped':
        return 'bg-purple-50 text-purple-700 border-purple-200';
      case 'completed':
        return 'bg-emerald-50 text-emerald-700 border-emerald-200';
      case 'cancelled':
        return 'bg-gray-50 text-gray-600 border-gray-200';
      default:
        return 'bg-amber-50 text-amber-800 border-amber-200';
    }
  };

  return (
    <div className="bg-primary pt-28 lg:pt-36 pb-24 min-h-screen">
      <div className="container-editorial max-w-4xl">
        <Link
          to="/shop"
          className="mb-8 inline-flex items-center gap-2 text-[11px] uppercase tracking-editorial text-dark/60 transition-colors hover:text-secondary"
        >
          <ArrowLeft size={14} /> Back to shop
        </Link>

        <h1 className="font-heading text-3xl sm:text-4xl text-dark mb-10">My Orders</h1>

        {loading ? (
          <div className="space-y-4">
            {[1, 2].map((i) => (
              <div key={i} className="animate-pulse bg-white/60 p-6 border border-dark/10 h-40" />
            ))}
          </div>
        ) : orders.length === 0 ? (
          <div className="bg-white/60 border border-dark/10 p-12 text-center">
            <Package size={40} className="mx-auto text-dark/20 mb-3 stroke-1" />
            <p className="font-heading text-xl text-dark">No orders yet</p>
            <p className="mt-2 text-xs font-light text-dark/60">
              When you place an order, its details and shipment tracking will appear here.
            </p>
            <Link
              to="/shop"
              className="mt-6 inline-flex items-center gap-2 bg-secondary px-6 py-3 text-[11px] uppercase tracking-editorial text-primary hover:bg-secondary-600 transition-colors"
            >
              Start Shopping
            </Link>
          </div>
        ) : (
          <div className="space-y-6">
            {orders.map((order) => (
              <div
                key={order.id}
                className="bg-white/70 border border-dark/10 p-6 transition-all hover:border-secondary/40"
              >
                {/* Order Top Bar */}
                <div className="flex flex-wrap items-center justify-between gap-4 border-b border-dark/10 pb-4">
                  <div>
                    <span className="text-[10px] uppercase tracking-editorial text-dark/50">
                      Order Reference
                    </span>
                    <p className="font-heading text-base text-dark">{order.order_number}</p>
                  </div>

                  <div className="flex items-center gap-4">
                    <div className="flex items-center gap-1.5 text-xs text-dark/60">
                      <Calendar size={14} />
                      <span>{new Date(order.created_at).toLocaleDateString('id-ID')}</span>
                    </div>

                    <span
                      className={`px-3 py-1 text-[10px] uppercase tracking-editorial font-medium border ${getStatusBadge(
                        order.status
                      )}`}
                    >
                      {order.status.replace('_', ' ')}
                    </span>
                  </div>
                </div>

                {/* Items List */}
                <div className="py-4 space-y-3 divide-y divide-dark/5">
                  {order.items.map((it) => (
                    <div key={it.id} className="pt-3 first:pt-0 flex items-center justify-between text-xs">
                      <div>
                        <p className="font-medium text-dark">{it.product_name}</p>
                        <p className="text-[10px] text-dark/50">
                          {it.color} / {it.size} · Qty: {it.quantity}
                        </p>
                      </div>
                      <span className="font-light text-dark">{formatIDR(it.subtotal)}</span>
                    </div>
                  ))}
                </div>

                {/* Bottom Courier & Total */}
                <div className="flex flex-wrap items-center justify-between gap-4 border-t border-dark/10 pt-4 text-xs">
                  <div className="flex items-center gap-2 text-dark/60">
                    <Truck size={15} className="text-secondary" />
                    <span>
                      {order.shipping_courier} ({order.shipping_service}) → {order.shipping_city}
                    </span>
                  </div>

                  <div className="flex items-baseline gap-2">
                    <span className="text-[10px] uppercase tracking-editorial text-dark/60">
                      Total Paid:
                    </span>
                    <span className="font-heading text-base text-dark">
                      {formatIDR(order.total)}
                    </span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
