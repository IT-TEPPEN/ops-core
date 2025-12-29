import type { Meta, StoryObj } from "@storybook/react";
import { CheckCircleIcon } from "./CheckCircleIcon";

const meta = {
  title: "UI/Icon/CheckCircleIcon",
  component: CheckCircleIcon,
  parameters: {
    layout: "centered",
    docs: {
      description: {
        component:
          "A check mark inside a circle icon. Used to indicate success, completion, or confirmation.",
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
} satisfies Meta<typeof CheckCircleIcon>;

export default meta;
type Story = StoryObj<typeof meta>;

/**
 * Default check circle with medium size (24px)
 */
export const Default: Story = {
  args: {},
};

/**
 * Small check circle (16px) - useful for inline success indicators
 */
export const Small: Story = {
  args: {
    size: 16,
  },
};

/**
 * Large check circle (32px) - useful for prominent success states
 */
export const Large: Story = {
  args: {
    size: 32,
  },
};

/**
 * Extra large check circle (48px) - useful for page-level success messages
 */
export const ExtraLarge: Story = {
  args: {
    size: 48,
  },
};

/**
 * Check circle with green color (success state)
 */
export const Green: Story = {
  args: {
    className: "text-green-600",
  },
};

/**
 * Check circle with blue color
 */
export const Blue: Story = {
  args: {
    className: "text-blue-600",
  },
};

/**
 * Check circle with custom styling
 */
export const WithBackground: Story = {
  render: () => (
    <div className="bg-green-50 dark:bg-green-900/20 p-4 rounded-lg border border-green-200 dark:border-green-800">
      <CheckCircleIcon className="text-green-600" size={32} />
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
        <CheckCircleIcon size={16} className="text-green-600" />
        <span className="text-xs text-gray-600">16px</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <CheckCircleIcon size={24} className="text-green-600" />
        <span className="text-xs text-gray-600">24px</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <CheckCircleIcon size={32} className="text-green-600" />
        <span className="text-xs text-gray-600">32px</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <CheckCircleIcon size={48} className="text-green-600" />
        <span className="text-xs text-gray-600">48px</span>
      </div>
    </div>
  ),
};

/**
 * Multiple colors comparison
 */
export const AllColors: Story = {
  render: () => (
    <div className="flex items-center gap-6">
      <div className="flex flex-col items-center gap-2">
        <CheckCircleIcon className="text-green-600" size={32} />
        <span className="text-xs text-gray-600">Green</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <CheckCircleIcon className="text-blue-600" size={32} />
        <span className="text-xs text-gray-600">Blue</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <CheckCircleIcon className="text-emerald-600" size={32} />
        <span className="text-xs text-gray-600">Emerald</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <CheckCircleIcon className="text-teal-600" size={32} />
        <span className="text-xs text-gray-600">Teal</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <CheckCircleIcon className="text-cyan-600" size={32} />
        <span className="text-xs text-gray-600">Cyan</span>
      </div>
    </div>
  ),
};

/**
 * Usage example in a success message
 */
export const InSuccessMessage: Story = {
  render: () => (
    <div className="bg-green-50 dark:bg-green-900/20 p-6 rounded-lg shadow border border-green-200 dark:border-green-800">
      <div className="flex items-center gap-3">
        <CheckCircleIcon className="text-green-600" size={24} />
        <div>
          <h3 className="font-semibold text-green-800 dark:text-green-200">
            Success!
          </h3>
          <p className="text-sm text-green-600 dark:text-green-400">
            Your changes have been saved successfully.
          </p>
        </div>
      </div>
    </div>
  ),
};
