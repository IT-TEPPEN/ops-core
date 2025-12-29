import type { Meta, StoryObj } from "@storybook/react";
import { SpinnerIcon } from "./SpinnerIcon";

const meta = {
  title: "UI/Icon/SpinnerIcon",
  component: SpinnerIcon,
  parameters: {
    layout: "centered",
    docs: {
      description: {
        component:
          "A spinning loading indicator icon. Automatically animates to indicate loading state.",
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
} satisfies Meta<typeof SpinnerIcon>;

export default meta;
type Story = StoryObj<typeof meta>;

/**
 * Default spinner with medium size (24px)
 */
export const Default: Story = {
  args: {},
};

/**
 * Small spinner (16px) - useful for inline loading indicators
 */
export const Small: Story = {
  args: {
    size: 16,
  },
};

/**
 * Large spinner (32px) - useful for page-level loading indicators
 */
export const Large: Story = {
  args: {
    size: 32,
  },
};

/**
 * Extra large spinner (48px) - useful for prominent loading states
 */
export const ExtraLarge: Story = {
  args: {
    size: 48,
  },
};

/**
 * Spinner with custom color (blue)
 */
export const Blue: Story = {
  args: {
    className: "text-blue-600",
  },
};

/**
 * Spinner with custom color (green)
 */
export const Green: Story = {
  args: {
    className: "text-green-600",
  },
};

/**
 * Spinner with custom color (red)
 */
export const Red: Story = {
  args: {
    className: "text-red-600",
  },
};

/**
 * Multiple sizes comparison
 */
export const AllSizes: Story = {
  render: () => (
    <div className="flex items-center gap-6">
      <div className="flex flex-col items-center gap-2">
        <SpinnerIcon size={16} />
        <span className="text-xs text-gray-600">16px</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <SpinnerIcon size={24} />
        <span className="text-xs text-gray-600">24px</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <SpinnerIcon size={32} />
        <span className="text-xs text-gray-600">32px</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <SpinnerIcon size={48} />
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
        <SpinnerIcon className="text-blue-600" size={32} />
        <span className="text-xs text-gray-600">Blue</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <SpinnerIcon className="text-green-600" size={32} />
        <span className="text-xs text-gray-600">Green</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <SpinnerIcon className="text-red-600" size={32} />
        <span className="text-xs text-gray-600">Red</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <SpinnerIcon className="text-yellow-600" size={32} />
        <span className="text-xs text-gray-600">Yellow</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <SpinnerIcon className="text-purple-600" size={32} />
        <span className="text-xs text-gray-600">Purple</span>
      </div>
    </div>
  ),
};
