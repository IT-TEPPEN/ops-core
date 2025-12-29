import type { Meta, StoryObj } from "@storybook/react";
import { LightningIcon } from "./LightningIcon";

const meta = {
  title: "UI/Icon/LightningIcon",
  component: LightningIcon,
  parameters: {
    layout: "centered",
    docs: {
      description: {
        component:
          "A lightning bolt icon. Used to indicate power, speed, connection, or quick actions.",
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
} satisfies Meta<typeof LightningIcon>;

export default meta;
type Story = StoryObj<typeof meta>;

/**
 * Default lightning icon with medium size (24px)
 */
export const Default: Story = {
  args: {},
};

/**
 * Small lightning icon (16px)
 */
export const Small: Story = {
  args: {
    size: 16,
  },
};

/**
 * Large lightning icon (32px)
 */
export const Large: Story = {
  args: {
    size: 32,
  },
};

/**
 * Lightning icon with yellow color
 */
export const Yellow: Story = {
  args: {
    className: "text-yellow-500",
  },
};

/**
 * Lightning icon with blue color
 */
export const Blue: Story = {
  args: {
    className: "text-blue-500",
  },
};

/**
 * Usage example in a connect button
 */
export const InConnectButton: Story = {
  render: () => (
    <button className="flex items-center gap-2 px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors">
      <LightningIcon size={20} />
      <span>Connect with GitHub</span>
    </button>
  ),
};

/**
 * Multiple sizes comparison
 */
export const AllSizes: Story = {
  render: () => (
    <div className="flex items-center gap-6">
      <div className="flex flex-col items-center gap-2">
        <LightningIcon size={16} className="text-yellow-500" />
        <span className="text-xs text-gray-600">16px</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <LightningIcon size={20} className="text-yellow-500" />
        <span className="text-xs text-gray-600">20px</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <LightningIcon size={24} className="text-yellow-500" />
        <span className="text-xs text-gray-600">24px</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <LightningIcon size={32} className="text-yellow-500" />
        <span className="text-xs text-gray-600">32px</span>
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
        <LightningIcon className="text-yellow-500" size={32} />
        <span className="text-xs text-gray-600">Yellow</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <LightningIcon className="text-blue-500" size={32} />
        <span className="text-xs text-gray-600">Blue</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <LightningIcon className="text-purple-500" size={32} />
        <span className="text-xs text-gray-600">Purple</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <LightningIcon className="text-green-500" size={32} />
        <span className="text-xs text-gray-600">Green</span>
      </div>
    </div>
  ),
};
