interface SubmitProps
  extends Omit<
    React.ButtonHTMLAttributes<HTMLButtonElement>,
    "type" | "disabled" | "className"
  > {
  isSubmitting: boolean;
  label: string;
}

export function UI_Form_Submit(props: SubmitProps) {
  return (
    <button
      type="submit"
      disabled={props.isSubmitting}
      className={`px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 ${
        props.isSubmitting ? "opacity-50 cursor-not-allowed" : ""
      }`}
    >
      {props.isSubmitting ? "Waiting..." : props.label}
    </button>
  );
}
