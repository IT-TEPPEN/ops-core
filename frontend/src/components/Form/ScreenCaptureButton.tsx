/**
 * ScreenCaptureButton component for capturing screen and uploading
 */

import { useState } from "react";
import { useScreenCapture } from "../../hooks/useScreenCapture";

export interface ScreenCaptureButtonProps {
  /** Execution record ID */
  executionRecordId: string;
  /** Execution step ID */
  executionStepId: string;
  /** Callback when capture is successful */
  onCaptureSuccess?: (attachment: AttachmentResponse) => void;
  /** Callback when capture fails */
  onCaptureError?: (error: Error) => void;
  /** Button variant */
  variant?: "primary" | "secondary" | "outline";
  /** Button size */
  size?: "sm" | "md" | "lg";
  /** Additional class name */
  className?: string;
}

export interface AttachmentResponse {
  id: string;
  execution_record_id: string;
  execution_step_id: string;
  file_name: string;
  file_size: number;
  mime_type: string;
  storage_type: string;
  uploaded_by: string;
  uploaded_at: string;
}

const variantClasses = {
  primary: "bg-blue-600 hover:bg-blue-700 text-white",
  secondary: "bg-gray-600 hover:bg-gray-700 text-white",
  outline: "border-2 border-gray-300 dark:border-gray-600 hover:border-gray-400 dark:hover:border-gray-500 text-gray-700 dark:text-gray-300",
};

const sizeClasses = {
  sm: "px-3 py-1.5 text-sm",
  md: "px-4 py-2 text-base",
  lg: "px-6 py-3 text-lg",
};

/**
 * Button component for capturing screen and uploading as attachment
 */
export function ScreenCaptureButton({
  executionRecordId,
  executionStepId,
  onCaptureSuccess,
  onCaptureError,
  variant = "primary",
  size = "md",
  className = "",
}: ScreenCaptureButtonProps) {
  const [capturing, setCapturing] = useState(false);
  const { captureScreen } = useScreenCapture();

  const handleCapture = async () => {
    setCapturing(true);
    try {
      const blob = await captureScreen();
      
      // Generate filename with timestamp
      const timestamp = new Date().toISOString().replace(/[:.]/g, "-");
      const filename = `screenshot-${timestamp}.png`;

      // Create FormData and upload
      const formData = new FormData();
      formData.append("file", blob, filename);
      formData.append("execution_step_id", executionStepId);

      const response = await fetch(
        `/api/v1/execution-records/${executionRecordId}/attachments`,
        {
          method: "POST",
          body: formData,
        }
      );

      if (!response.ok) {
        throw new Error(`Upload failed: ${response.statusText}`);
      }

      const attachment: AttachmentResponse = await response.json();
      onCaptureSuccess?.(attachment);
    } catch (error) {
      onCaptureError?.(error as Error);
    } finally {
      setCapturing(false);
    }
  };

  return (
    <button
      onClick={handleCapture}
      disabled={capturing}
      className={`
        inline-flex items-center justify-center
        rounded-lg font-medium
        transition-colors
        disabled:opacity-50 disabled:cursor-not-allowed
        ${variantClasses[variant]}
        ${sizeClasses[size]}
        ${className}
      `}
    >
      {capturing ? (
        <>
          <svg
            className="animate-spin -ml-1 mr-2 h-5 w-5"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              className="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              strokeWidth="4"
            />
            <path
              className="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            />
          </svg>
          Capturing...
        </>
      ) : (
        <>
          <svg
            className="w-5 h-5 mr-2"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z"
            />
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M15 13a3 3 0 11-6 0 3 3 0 016 0z"
            />
          </svg>
          Capture Screen
        </>
      )}
    </button>
  );
}
