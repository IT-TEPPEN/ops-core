import type { Meta, StoryObj } from "@storybook/react";
import { CheckCircleFilledIcon } from "./CheckCircleFilledIcon";

const meta = {
  title: "UI/Icon/CheckCircleFilledIcon",
  component: CheckCircleFilledIcon,
  parameters: {
    layout: "centered",
    docs: {
      description: {
        component:
          "A filled check circle icon. Used to indicate selected items, completed states, or confirmed choices.",
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
} satisfies Meta<typeof CheckCircleFilledIcon>;

export default meta;
type Story = StoryObj<typeof meta>;

/**
 * Default filled check circle with medium size (24px)
 */
export const Default: Story = {
  args: {},
};

/**
 * Small filled check circle (16px)
 */
export const Small: Story = {
  args: {
    size: 16,
  },
};

/**
 * Large filled check circle (32px)
 */
export const Large: Story = {
  args: {
    size: 32,
  },
};

/**
 * Filled check circle with blue color (selected state)
 */
export const Blue: Story = {
  args: {
    className: "text-blue-600",
  },
};

/**
 * Filled check circle with green color
 */
export const Green: Story = {
  args: {
    className: "text-green-600",
  },
};

/**
 * Usage example in a selectable list
 */
export const InSelectableList: Story = {
  render: () => (
    <div className="space-y-2">
      <div className="flex items-center gap-3 p-3 border border-gray-200 rounded-lg">
        <div className="flex-1">
          <p className="font-medium">Option 1</p>
          <p className="text-sm text-gray-600">This option is selected</p>
        </div>
        <CheckCircleFilledIcon className="text-blue-600 shrink-0" size={20} />
      </div>
      <div className="flex items-center gap-3 p-3 border border-gray-200 rounded-lg opacity-50">
        <div className="flex-1">
          <p className="font-medium">Option 2</p>
          <p className="text-sm text-gray-600">This option is not selected</p>
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
        <CheckCircleFilledIcon size={16} className="text-blue-600" />
        <span className="text-xs text-gray-600">16px</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <CheckCircleFilledIcon size={20} className="text-blue-600" />
        <span className="text-xs text-gray-600">20px</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <CheckCircleFilledIcon size={24} className="text-blue-600" />
        <span className="text-xs text-gray-600">24px</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <CheckCircleFilledIcon size={32} className="text-blue-600" />
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
        <CheckCircleFilledIcon className="text-blue-600" size={32} />
        <span className="text-xs text-gray-600">Blue</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <CheckCircleFilledIcon className="text-green-600" size={32} />
        <span className="text-xs text-gray-600">Green</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <CheckCircleFilledIcon className="text-purple-600" size={32} />
        <span className="text-xs text-gray-600">Purple</span>
      </div>
      <div className="flex flex-col items-center gap-2">
        <CheckCircleFilledIcon className="text-indigo-600" size={32} />
        <span className="text-xs text-gray-600">Indigo</span>
      </div>
    </div>
  ),
};
