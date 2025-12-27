import { InputHTMLAttributes } from "react";

export function UI_Form_Input(props: InputHTMLAttributes<HTMLInputElement>) {
  const { className, ...others } = props;

  return (
    <input
      className={`w-full px-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white ${
        className ?? ""
      }`}
      {...others}
    />
  );
}
