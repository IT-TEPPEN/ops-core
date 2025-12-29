import type { Meta, StoryObj } from "@storybook/react";
import { XCircleIcon } from "./XCircleIcon";

const meta = {
  title: "UI/Icon/XCircleIcon",
  component: XCircleIcon,
  parameters: {
    layout: "centered",
    docs: {
      description: {
        component:
          "An X mark inside a circle icon. Used to indicate errors, failures, or invalid states.",
      },
    },
  },
  tags: ["autodocs"],
  argTypes: {
    size: {
      control: { type: "number", min: 16, max: 96, step: 8 },
      description: "Icon size in pixels",
    },
    className: {
      control: "text",
      description: "Additional CSS classes for styling",
    },
  },
} satisfies Meta<typeof XCircleIcon>;

export default meta;
type Story = StoryObj<typeof meta>;

/**
 * Default X circle with medium size (24px)
 */
export const Default: Story = {
  args: {},
};

/**
 * Small X circle (16px)
 */
export const Small: Story = {
  args: {
    size: 16,
  },
};

/**
 * Large X circle (32px)
 */
export const Large: Story = {
  args: {
    size: 32,
  },
};

/**
 * X circle with red color (error state)
 */
export const Red: Story = {
  args: {
    className: "text-red-400",
  },
};

/**
 * Usage example in an error message
 */
export const InErrorMessage: Story = {
  render: () => (
    <div className="bg-red-50 dark:bg-red-900/20 p-4 rounded-lg border border-red-200 dark:border-red-800">
      <div className="flex">
        <div className="shrink-0">
          <XCircleIcon className="text-red-400" size={20} />
        </div>
        <div className="ml-3">
          <h3 className="text-sm font-medium text-red-800 dark:text-red-200">
            Error
          </h3>
          <p className="text-sm text-red-700 dark:text-red-300 mt-1">
            Something went wrong. Please try again.
          </p>
        </div>
      </div>
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
        <XCircleIcon size={16} className="text-red-400" />
        <span className="text-xs text-gray-600">16px</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <XCircleIcon size={20} className="text-red-400" />
        <span className="text-xs text-gray-600">20px</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <XCircleIcon size={24} className="text-red-400" />
        <span className="text-xs text-gray-600">24px</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <XCircleIcon size={32} className="text-red-400" />
        <span className="text-xs text-gray-600">32px</span>
      </div>
    </div>
  ),
};
