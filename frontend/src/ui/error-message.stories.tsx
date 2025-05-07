import type { Meta, StoryObj } from "@storybook/react";

import ErrorMessage from "./ErrorMessage";

export const ActionsData = {};

const meta = {
  component: ErrorMessage,
  title: "ErrorMessage",
  tags: ["autodocs"],
  //👇 Our exports that end in "Data" are not stories.
  excludeStories: /.*Data$/,
  args: {
    ...ActionsData,
  },
} satisfies Meta<typeof ErrorMessage>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    message: "This is an error message",
  },
  render: (args) => <ErrorMessage {...args} />,
};
