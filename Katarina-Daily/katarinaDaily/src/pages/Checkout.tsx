import React, { useEffect, useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { ArrowLeft, CheckCircle2, ShieldCheck, Truck, CreditCard, ShoppingBag, AlertCircle, ExternalLink } from 'lucide-react';
import { useCart } from '@/context/CartContext';
import { useAuth } from '@/context/AuthContext';
import { formatIDR } from '@/services/catalogService';
import {
  getAddresses,
  createAddress,
  getShippingRates,
  checkout as apiCheckout,
  type Address,
  type ShippingRateOption,
  type Order,
} from '@/services/orderService';

export default function Checkout() {
  const navigate = useNavigate();
  const { cart, clear } = useCart();
  const { isAuthenticated, openAuthModal } = useAuth();

  const [addresses, setAddresses] = useState<Address[]>([]);
  const [selectedAddressId, setSelectedAddressId] = useState<string>('');
  const [shippingRates, setShippingRates] = useState<ShippingRateOption[]>([]);
  const [selectedCourier, setSelectedCourier] = useState<string>('JNE');
  const [selectedService, setSelectedService] = useState<string>('REG');
  const [shippingCost, setShippingCost] = useState<number>(18000);
  const [voucherCode, setVoucherCode] = useState('');
  const [loadingRates, setLoadingRates] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [orderComplete, setOrderComplete] = useState<Order | null>(null);
  const [checkoutError, setCheckoutError] = useState<string | null>(null);

  // Address form fields
  const [recipientName, setRecipientName] = useState('');
  const [phone, setPhone] = useState('');
  const [addressLine, setAddressLine] = useState('');
  const [city, setCity] = useState('Jakarta Selatan');
  const [province, setProvince] = useState('DKI Jakarta');
  const [postalCode, setPostalCode] = useState('12190');

  useEffect(() => {
    if (isAuthenticated) {
      getAddresses().then((list) => {
        setAddresses(list);
        if (list.length > 0) {
          const defaultAddr = list.find((a) => a.is_default) || list[0];
          setSelectedAddressId(defaultAddr.id);
          setRecipientName(defaultAddr.recipient_name);
          setPhone(defaultAddr.phone);
          setAddressLine(defaultAddr.address);
          setCity(defaultAddr.city);
          setProvince(defaultAddr.province);
          setPostalCode(defaultAddr.postal_code);
        }
      });
    }
  }, [isAuthenticated]);

  useEffect(() => {
    async function loadRates() {
      setLoadingRates(true);
      const rates = await getShippingRates(selectedAddressId || 'addr-temp', selectedCourier);
      setShippingRates(rates);
      if (rates.length > 0) {
        const first = rates[0];
        setSelectedCourier(first.courier);
        setSelectedService(first.service);
        setShippingCost(typeof first.cost === 'string' ? parseFloat(first.cost) : first.cost);
      }
      setLoadingRates(false);
    }
    loadRates();
  }, [selectedAddressId, selectedCourier]);

  const handleCourierChange = (opt: ShippingRateOption) => {
    setSelectedCourier(opt.courier);
    setSelectedService(opt.service);
    setShippingCost(typeof opt.cost === 'string' ? parseFloat(opt.cost) : opt.cost);
  };

  const handlePlaceOrder = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!isAuthenticated) {
      openAuthModal('login');
      return;
    }

    setSubmitting(true);
    setCheckoutError(null);
    try {
      let finalAddressId = selectedAddressId;

      // If no saved address chosen, create one first
      if (!finalAddressId) {
        const newAddr = await createAddress({
          recipient_name: recipientName,
          phone,
          address: addressLine,
          city,
          province,
          postal_code: postalCode,
          is_default: true,
        });
        finalAddressId = newAddr.id;
      }

      const order = await apiCheckout({
        address_id: finalAddressId,
        shipping_courier: selectedCourier,
        shipping_service: selectedService,
        voucher_code: voucherCode.trim() || undefined,
      });

      await clear();
      setOrderComplete(order);
    } catch (err: any) {
      console.error('Checkout error:', err);
      setCheckoutError(
        err?.message || 'Gagal memproses pesanan. Silakan periksa kembali data Anda dan coba lagi.'
      );
    } finally {
      setSubmitting(false);
    }
  };

  if (orderComplete) {
    const isPendingPayment = orderComplete.status.toLowerCase().includes('pending');
    const paymentUrl = orderComplete.payment?.redirect_url;

    return (
      <div className="bg-primary pt-32 lg:pt-40 pb-24 min-h-screen">
        <div className="container-editorial max-w-xl text-center">
          <div className="mx-auto mb-6 flex h-16 w-16 items-center justify-center rounded-full bg-secondary/15 text-secondary">
            <CheckCircle2 size={36} />
          </div>
          <p className="text-[11px] uppercase tracking-editorial text-secondary mb-2">
            {isPendingPayment ? 'Order Created · Awaiting Payment' : 'Order Confirmed'}
          </p>
          <h1 className="font-heading text-3xl sm:text-4xl text-dark">
            {isPendingPayment ? 'Selesaikan Pembayaran Anda' : 'Terima Kasih Atas Pesanan Anda!'}
          </h1>
          <p className="mt-3 text-xs font-light text-dark/70">
            Nomor Pesanan: <strong className="font-medium text-dark">{orderComplete.order_number}</strong>
          </p>

          {isPendingPayment && paymentUrl && (
            <div className="mt-8 bg-amber-50 border border-amber-200 p-6 text-left">
              <p className="text-xs font-medium text-amber-900 mb-1">
                Langkah Pembayaran Midtrans
              </p>
              <p className="text-[11px] text-amber-800 font-light mb-4">
                Pesanan Anda telah berhasil dibuat. Silakan klik tombol di bawah untuk melanjutkan pembayaran via Midtrans (GoPay, QRIS, Bank Transfer, atau Kartu Kredit).
              </p>
              <a
                href={paymentUrl}
                target="_blank"
                rel="noopener noreferrer"
                className="w-full inline-flex items-center justify-center gap-2 bg-secondary py-3.5 px-6 text-xs uppercase tracking-editorial text-primary hover:bg-secondary-600 transition-colors font-medium shadow-sm"
              >
                <CreditCard size={15} /> Bayar Sekarang (Midtrans) <ExternalLink size={13} />
              </a>
            </div>
          )}

          <div className="mt-8 bg-white/70 border border-dark/10 p-6 text-left space-y-4">
            <div className="flex justify-between text-xs border-b border-dark/5 pb-3">
              <span className="text-dark/60">Penerima</span>
              <span className="font-medium text-dark">{orderComplete.shipping_recipient_name}</span>
            </div>
            <div className="flex justify-between text-xs border-b border-dark/5 pb-3">
              <span className="text-dark/60">Alamat Pengiriman</span>
              <span className="font-medium text-dark text-right max-w-xs">
                {orderComplete.shipping_address}, {orderComplete.shipping_city} ({orderComplete.shipping_postal_code})
              </span>
            </div>
            <div className="flex justify-between text-xs border-b border-dark/5 pb-3">
              <span className="text-dark/60">Kurir</span>
              <span className="font-medium text-dark">
                {orderComplete.shipping_courier} - {orderComplete.shipping_service}
              </span>
            </div>
            <div className="flex justify-between text-xs pt-2">
              <span className="text-xs uppercase tracking-editorial font-medium text-dark">Total Tagihan</span>
              <span className="font-heading text-base text-dark">
                {formatIDR(orderComplete.total)}
              </span>
            </div>
          </div>

          <div className="mt-8 flex flex-col sm:flex-row items-center justify-center gap-4">
            <Link
              to="/orders"
              className="w-full sm:w-auto bg-secondary px-8 py-3.5 text-xs uppercase tracking-editorial text-primary hover:bg-secondary-600 transition-colors"
            >
              Lihat Riwayat Pesanan
            </Link>
            <Link
              to="/shop"
              className="w-full sm:w-auto border border-dark/20 px-8 py-3.5 text-xs uppercase tracking-editorial text-dark hover:border-secondary hover:text-secondary transition-colors"
            >
              Lanjut Belanja
            </Link>
          </div>
        </div>
      </div>
    );
  }

  if (cart.items.length === 0) {
    return (
      <div className="bg-primary pt-32 lg:pt-40 pb-24 min-h-screen">
        <div className="container-editorial max-w-md text-center">
          <ShoppingBag size={44} className="mx-auto text-dark/30 mb-4 stroke-1" />
          <h2 className="font-heading text-2xl text-dark">Your bag is empty</h2>
          <p className="mt-2 text-xs font-light text-dark/60">
            Please add items to your shopping bag before proceeding to checkout.
          </p>
          <Link
            to="/shop"
            className="mt-6 inline-flex items-center gap-2 bg-secondary px-8 py-3.5 text-xs uppercase tracking-editorial text-primary hover:bg-secondary-600 transition-colors"
          >
            Explore Collection
          </Link>
        </div>
      </div>
    );
  }

  const grandTotal = cart.subtotal + shippingCost;

  return (
    <div className="bg-primary pt-28 lg:pt-36 pb-24 lg:pb-32 min-h-screen">
      <div className="container-editorial max-w-5xl">
        <Link
          to="/shop"
          className="mb-8 inline-flex items-center gap-2 text-[11px] uppercase tracking-editorial text-dark/60 transition-colors hover:text-secondary"
        >
          <ArrowLeft size={14} /> Back to shopping
        </Link>

        <h1 className="font-heading text-3xl sm:text-4xl text-dark mb-10">Checkout</h1>

        {!isAuthenticated && (
          <div className="mb-8 bg-secondary/10 border border-secondary/20 p-4 flex items-center justify-between">
            <p className="text-xs text-dark/80">
              Have an account? Sign in for saved addresses and live tracking.
            </p>
            <button
              onClick={() => openAuthModal('login')}
              className="bg-secondary px-4 py-2 text-[10px] uppercase tracking-editorial text-primary hover:bg-secondary-600 transition-colors"
            >
              Sign In
            </button>
          </div>
        )}

        <form onSubmit={handlePlaceOrder} className="grid grid-cols-1 lg:grid-cols-12 gap-12">
          {/* Left: Shipping & Details */}
          <div className="lg:col-span-7 space-y-8">
            {/* Delivery Address Section */}
            <div className="bg-white/70 border border-dark/10 p-6 space-y-4">
              <div className="flex items-center gap-2 border-b border-dark/10 pb-3">
                <Truck size={18} className="text-secondary" />
                <h3 className="font-heading text-lg text-dark">Shipping Destination</h3>
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label className="block text-[10px] uppercase tracking-editorial text-dark/70 mb-1.5">
                    Recipient Name
                  </label>
                  <input
                    type="text"
                    required
                    value={recipientName}
                    onChange={(e) => setRecipientName(e.target.value)}
                    placeholder="Jane Doe"
                    className="w-full bg-white border border-dark/15 px-3.5 py-2.5 text-xs text-dark focus:border-secondary focus:outline-none"
                  />
                </div>
                <div>
                  <label className="block text-[10px] uppercase tracking-editorial text-dark/70 mb-1.5">
                    Phone Number
                  </label>
                  <input
                    type="tel"
                    required
                    value={phone}
                    onChange={(e) => setPhone(e.target.value)}
                    placeholder="081234567890"
                    className="w-full bg-white border border-dark/15 px-3.5 py-2.5 text-xs text-dark focus:border-secondary focus:outline-none"
                  />
                </div>
              </div>

              <div>
                <label className="block text-[10px] uppercase tracking-editorial text-dark/70 mb-1.5">
                  Street Address
                </label>
                <textarea
                  required
                  rows={2}
                  value={addressLine}
                  onChange={(e) => setAddressLine(e.target.value)}
                  placeholder="Street name, house/apartment number"
                  className="w-full bg-white border border-dark/15 px-3.5 py-2 text-xs text-dark focus:border-secondary focus:outline-none"
                />
              </div>

              <div className="grid grid-cols-3 gap-3">
                <div>
                  <label className="block text-[10px] uppercase tracking-editorial text-dark/70 mb-1.5">
                    City
                  </label>
                  <input
                    type="text"
                    required
                    value={city}
                    onChange={(e) => setCity(e.target.value)}
                    placeholder="Jakarta Selatan"
                    className="w-full bg-white border border-dark/15 px-3 py-2 text-xs text-dark focus:border-secondary focus:outline-none"
                  />
                </div>
                <div>
                  <label className="block text-[10px] uppercase tracking-editorial text-dark/70 mb-1.5">
                    Province
                  </label>
                  <input
                    type="text"
                    required
                    value={province}
                    onChange={(e) => setProvince(e.target.value)}
                    placeholder="DKI Jakarta"
                    className="w-full bg-white border border-dark/15 px-3 py-2 text-xs text-dark focus:border-secondary focus:outline-none"
                  />
                </div>
                <div>
                  <label className="block text-[10px] uppercase tracking-editorial text-dark/70 mb-1.5">
                    Postal Code
                  </label>
                  <input
                    type="text"
                    required
                    value={postalCode}
                    onChange={(e) => setPostalCode(e.target.value)}
                    placeholder="12190"
                    className="w-full bg-white border border-dark/15 px-3 py-2 text-xs text-dark focus:border-secondary focus:outline-none"
                  />
                </div>
              </div>
            </div>

            {/* Courier Selection */}
            <div className="bg-white/70 border border-dark/10 p-6 space-y-4">
              <div className="flex items-center justify-between border-b border-dark/10 pb-3">
                <h3 className="font-heading text-lg text-dark">Courier & Delivery Service</h3>
                {loadingRates && (
                  <span className="text-[10px] text-dark/50">Fetching rates...</span>
                )}
              </div>

              <div className="space-y-2.5">
                {shippingRates.map((rate, idx) => (
                  <label
                    key={`${rate.courier}-${rate.service}-${idx}`}
                    className={`flex items-center justify-between p-3.5 border cursor-pointer transition-colors ${
                      selectedCourier === rate.courier && selectedService === rate.service
                        ? 'border-secondary bg-secondary/5'
                        : 'border-dark/10 hover:border-dark/30'
                    }`}
                  >
                    <div className="flex items-center gap-3">
                      <input
                        type="radio"
                        name="courierOption"
                        checked={selectedCourier === rate.courier && selectedService === rate.service}
                        onChange={() => handleCourierChange(rate)}
                        className="accent-secondary"
                      />
                      <div>
                        <p className="text-xs font-medium text-dark">
                          {rate.courier} · {rate.service}
                        </p>
                        <p className="text-[10px] text-dark/60">{rate.description} ({rate.etd})</p>
                      </div>
                    </div>
                    <span className="text-xs font-medium text-dark">{formatIDR(rate.cost)}</span>
                  </label>
                ))}
              </div>
            </div>

            {/* Payment Method Note */}
            <div className="bg-white/70 border border-dark/10 p-6 flex items-start gap-3.5">
              <CreditCard size={20} className="text-secondary shrink-0 mt-0.5" />
              <div>
                <h4 className="text-xs font-medium uppercase tracking-editorial text-dark">
                  Integrated Midtrans Payment Gateway
                </h4>
                <p className="text-xs font-light text-dark/60 mt-1">
                  Supports QRIS, BCA/Mandiri/BNI Virtual Accounts, GoPay, and Credit Cards with instant automated verification.
                </p>
              </div>
            </div>
          </div>

          {/* Right: Order Summary */}
          <div className="lg:col-span-5">
            <div className="bg-white/70 border border-dark/10 p-6 space-y-6 sticky top-28">
              <h3 className="font-heading text-lg text-dark border-b border-dark/10 pb-3">
                Order Summary ({cart.total_items} items)
              </h3>

              {/* Items List */}
              <div className="max-h-64 overflow-y-auto space-y-3.5 divide-y divide-dark/5 pr-1">
                {cart.items.map((item) => (
                  <div key={item.id} className="pt-3.5 first:pt-0 flex gap-3 items-center">
                    <img
                      src={
                        item.product_image ||
                        'https://images.pexels.com/photos/22441297/pexels-photo-22441297.jpeg?auto=compress&cs=tinysrgb&h=900'
                      }
                      alt={item.product_name}
                      className="h-14 w-12 object-cover bg-primary-100 border border-dark/5"
                    />
                    <div className="flex-1 min-w-0">
                      <p className="text-xs font-medium text-dark truncate">{item.product_name}</p>
                      <p className="text-[10px] text-dark/50">
                        {item.color} / {item.size} · Qty: {item.quantity}
                      </p>
                    </div>
                    <span className="text-xs font-light text-dark">{formatIDR(item.subtotal)}</span>
                  </div>
                ))}
              </div>

              {/* Voucher code input */}
              <div className="flex gap-2">
                <input
                  type="text"
                  value={voucherCode}
                  onChange={(e) => setVoucherCode(e.target.value.toUpperCase())}
                  placeholder="Voucher code (optional)"
                  className="flex-1 bg-white border border-dark/15 px-3 py-2 text-xs text-dark focus:border-secondary focus:outline-none"
                />
                <button
                  type="button"
                  onClick={() => alert(`Voucher "${voucherCode}" applied!`)}
                  className="border border-dark/20 px-3.5 py-2 text-[10px] uppercase tracking-editorial text-dark hover:border-secondary hover:text-secondary transition-colors"
                >
                  Apply
                </button>
              </div>

              {/* Price Breakdown */}
              <div className="space-y-2 border-t border-dark/10 pt-4 text-xs font-light">
                <div className="flex justify-between text-dark/70">
                  <span>Subtotal</span>
                  <span>{formatIDR(cart.subtotal)}</span>
                </div>
                <div className="flex justify-between text-dark/70">
                  <span>Shipping ({selectedCourier})</span>
                  <span>{formatIDR(shippingCost)}</span>
                </div>
                <div className="flex justify-between text-base font-normal text-dark pt-2 border-t border-dark/5">
                  <span className="font-heading">Total</span>
                  <span className="font-heading text-lg">{formatIDR(grandTotal)}</span>
                </div>
              </div>

              {checkoutError && (
                <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 text-xs flex items-center gap-2">
                  <AlertCircle size={16} className="shrink-0 text-red-500" />
                  <span>{checkoutError}</span>
                </div>
              )}

              <button
                type="submit"
                disabled={submitting}
                className="w-full bg-secondary py-4 text-xs uppercase tracking-editorial text-primary hover:bg-secondary-600 transition-all duration-300 disabled:opacity-50"
              >
                {submitting ? 'Placing Order...' : `Place Order · ${formatIDR(grandTotal)}`}
              </button>

              <div className="flex items-center justify-center gap-2 text-[10px] text-dark/50">
                <ShieldCheck size={14} className="text-secondary" />
                <span>SSL Encrypted · Secure Checkout</span>
              </div>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
}
