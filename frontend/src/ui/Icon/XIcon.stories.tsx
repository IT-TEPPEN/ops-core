import type { Meta, StoryObj } from "@storybook/react";
import { XIcon } from "./XIcon";

const meta = {
  title: "UI/Icon/XIcon",
  component: XIcon,
  parameters: {
    layout: "centered",
    docs: {
      description: {
        component:
          "A simple X mark icon. Commonly used for close buttons in dialogs, modals, and dismissible components.",
      },
    },
  },
  tags: ["autodocs"],
  argTypes: {
    size: {
      control: { type: "number", min: 12, max: 48, step: 4 },
      description: "Icon size in pixels",
    },
    className: {
      control: "text",
      description: "Additional CSS classes for styling",
    },
  },
} satisfies Meta<typeof XIcon>;

export default meta;
type Story = StoryObj<typeof meta>;

/**
 * Default X icon with medium size (24px)
 */
export const Default: Story = {
  args: {},
};

/**
 * Small X icon (16px) - common for small UI elements
 */
export const Small: Story = {
  args: {
    size: 16,
  },
};

/**
 * Large X icon (32px) - for prominent close buttons
 */
export const Large: Story = {
  args: {
    size: 32,
  },
};

/**
 * X icon with gray color
 */
export const Gray: Story = {
  args: {
    className: "text-gray-500 dark:text-gray-400",
  },
};

/**
 * Usage example in a modal header
 */
export const InModalHeader: Story = {
  render: () => (
    <div className="flex items-center justify-between p-4 border-b border-gray-200 dark:border-gray-700">
      <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
        Modal Title
      </h2>
      <button
        className="p-1 rounded hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
        aria-label="Close"
      >
        <XIcon className="text-gray-500 dark:text-gray-400" size={20} />
      </button>
    </div>
  ),
};

/**
 * Usage example in a toast notification
 */
export const InToast: Story = {
  render: () => (
    <div className="flex items-center gap-3 p-3 bg-gray-900 text-white rounded-lg shadow-lg">
      <p className="flex-1 text-sm">This is a notification message</p>
      <button className="shrink-0 hover:opacity-70" aria-label="Dismiss">
        <XIcon size={16} />
      </button>
    </div>
  ),
};

/**
 * Multiple sizes comparison
 */
export const AllSizes: Story = {
  render: () => (
    <div className="flex items-center gap-6">
      <div className="flex flex-col items-center gap-2">
        <XIcon size={12} className="text-gray-600" />
        <span className="text-xs text-gray-600">12px</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <XIcon size={16} className="text-gray-600" />
        <span className="text-xs text-gray-600">16px</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <XIcon size={20} className="text-gray-600" />
        <span className="text-xs text-gray-600">20px</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <XIcon size={24} className="text-gray-600" />
        <span className="text-xs text-gray-600">24px</span>
      </div>
    </div>
  ),
};
