/**
 * Toast notification component and context
 */

import { createContext, useContext, useState, useCallback, useEffect } from "react";
import type { ReactNode } from "react";
import type { Toast as ToastType, Variant } from "../../types/ui";

// Toast context
interface ToastContextValue {
  toasts: ToastType[];
  addToast: (type: Variant, message: string, duration?: number) => void;
  removeToast: (id: string) => void;
}

const ToastContext = createContext<ToastContextValue | null>(null);

export function useToast() {
  const context = useContext(ToastContext);
  if (!context) {
    throw new Error("useToast must be used within a ToastProvider");
  }
  return context;
}

// Toast provider
export interface ToastProviderProps {
  children: ReactNode;
  /** Default duration in milliseconds */
  defaultDuration?: number;
}

export function ToastProvider({
  children,
  defaultDuration = 5000,
}: ToastProviderProps) {
  const [toasts, setToasts] = useState<ToastType[]>([]);

  const addToast = useCallback(
    (type: Variant, message: string, duration = defaultDuration) => {
      const id = Math.random().toString(36).substr(2, 9);
      setToasts((prev) => [...prev, { id, type, message, duration }]);
    },
    [defaultDuration]
  );

  const removeToast = useCallback((id: string) => {
    setToasts((prev) => prev.filter((toast) => toast.id !== id));
  }, []);

  return (
    <ToastContext.Provider value={{ toasts, addToast, removeToast }}>
      {children}
      <ToastContainer toasts={toasts} onRemove={removeToast} />
    </ToastContext.Provider>
  );
}

// Toast container
interface ToastContainerProps {
  toasts: ToastType[];
  onRemove: (id: string) => void;
}

function ToastContainer({ toasts, onRemove }: ToastContainerProps) {
  if (toasts.length === 0) return null;

  return (
    <div className="fixed bottom-4 right-4 z-50 flex flex-col gap-2">
      {toasts.map((toast) => (
        <ToastItem key={toast.id} toast={toast} onRemove={onRemove} />
      ))}
    </div>
  );
}

// Single toast item
interface ToastItemProps {
  toast: ToastType;
  onRemove: (id: string) => void;
}

const variantStyles: Record<Variant, { bg: string; icon: string }> = {
  primary: {
    bg: "bg-blue-500",
    icon: "ℹ️",
  },
  secondary: {
    bg: "bg-gray-500",
    icon: "📝",
  },
  success: {
    bg: "bg-green-500",
    icon: "✓",
  },
  warning: {
    bg: "bg-yellow-500",
    icon: "⚠",
  },
  danger: {
    bg: "bg-red-500",
    icon: "✕",
  },
  info: {
    bg: "bg-cyan-500",
    icon: "ℹ",
  },
};

function ToastItem({ toast, onRemove }: ToastItemProps) {
  useEffect(() => {
    if (toast.duration && toast.duration > 0) {
      const timer = setTimeout(() => {
        onRemove(toast.id);
      }, toast.duration);
      return () => clearTimeout(timer);
    }
  }, [toast.id, toast.duration, onRemove]);

  const styles = variantStyles[toast.type];

  return (
    <div
      className={`
        ${styles.bg} text-white
        px-4 py-3 rounded-lg shadow-lg
        flex items-center gap-3
        min-w-[200px] max-w-md
        animate-in slide-in-from-right
      `}
      role="alert"
    >
      <span className="flex-shrink-0">{styles.icon}</span>
      <p className="flex-1 text-sm">{toast.message}</p>
      <button
        onClick={() => onRemove(toast.id)}
        className="flex-shrink-0 hover:opacity-70"
        aria-label="Dismiss"
      >
        <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
          <path
            fillRule="evenodd"
            d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
            clipRule="evenodd"
          />
        </svg>
      </button>
    </div>
  );
}

// Standalone toast component for simple usage
export interface ToastProps {
  /** Toast type */
  type?: Variant;
  /** Toast message */
  message: string;
  /** Whether the toast is visible */
  isVisible: boolean;
  /** Close callback */
  onClose: () => void;
  /** Duration in milliseconds */
  duration?: number;
}

export function Toast({
  type = "info",
  message,
  isVisible,
  onClose,
  duration = 5000,
}: ToastProps) {
  useEffect(() => {
    if (isVisible && duration > 0) {
      const timer = setTimeout(onClose, duration);
      return () => clearTimeout(timer);
    }
  }, [isVisible, duration, onClose]);

  if (!isVisible) return null;

  const styles = variantStyles[type];

  return (
    <div className="fixed bottom-4 right-4 z-50">
      <div
        className={`
          ${styles.bg} text-white
          px-4 py-3 rounded-lg shadow-lg
          flex items-center gap-3
          min-w-[200px] max-w-md
        `}
        role="alert"
      >
        <span className="flex-shrink-0">{styles.icon}</span>
        <p className="flex-1 text-sm">{message}</p>
        <button
          onClick={onClose}
          className="flex-shrink-0 hover:opacity-70"
          aria-label="Dismiss"
        >
          <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
            <path
              fillRule="evenodd"
              d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
              clipRule="evenodd"
            />
          </svg>
        </button>
      </div>
    </div>
  );
}
