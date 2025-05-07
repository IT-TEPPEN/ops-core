import { defineWorkspace } from "vitest/config";
import { storybookTest } from "@storybook/experimental-addon-test/vitest-plugin";

// More info at: https://storybook.js.org/docs/writing-tests/vitest-plugin
export default defineWorkspace([
  "vite.config.ts",
  {
    extends: "vite.config.ts",
    plugins: [
      // See options at: https://storybook.js.org/docs/writing-tests/vitest-plugin#storybooktest
      storybookTest(),
    ],
    test: {
      name: "storybook",
      browser: {
        // vitest のブラウザモードを有効にする
        enabled: true,
        provider: "playwright",
        // Using instances array with proper object format instead of name field (Vitest 3 approach)
        instances: [
          {
            name: "chromium",
            browser: "chromium",
            headless: true,
          },
        ],
      },
      setupFiles: ["./.storybook/vitest.setup.ts"],
    },
  },
]);
