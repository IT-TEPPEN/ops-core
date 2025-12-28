/**
 * AttachmentUploader component for uploading files
 */

import { useState, useRef, type ChangeEvent } from "react";

export interface AttachmentUploaderProps {
  /** Execution record ID */
  executionRecordId: string;
  /** Execution step ID */
  executionStepId: string;
  /** Callback when upload is successful */
  onUploadSuccess?: (attachment: AttachmentResponse) => void;
  /** Callback when upload fails */
  onUploadError?: (error: Error) => void;
  /** Accepted file types */
  accept?: string;
  /** Maximum file size in bytes */
  maxSize?: number;
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

/**
 * Component for uploading attachment files
 */
export function AttachmentUploader({
  executionRecordId,
  executionStepId,
  onUploadSuccess,
  onUploadError,
  accept = "image/*,.pdf,.doc,.docx,.txt,.log",
  maxSize = 10 * 1024 * 1024, // 10MB default
  className = "",
}: AttachmentUploaderProps) {
  const [uploading, setUploading] = useState(false);
  const [dragOver, setDragOver] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const uploadFile = async (file: File) => {
    if (file.size > maxSize) {
      const error = new Error(`File size exceeds ${Math.round(maxSize / 1024 / 1024)}MB limit`);
      onUploadError?.(error);
      return;
    }

    setUploading(true);
    try {
      const formData = new FormData();
      formData.append("file", file);
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
      onUploadSuccess?.(attachment);
    } catch (error) {
      onUploadError?.(error as Error);
    } finally {
      setUploading(false);
    }
  };

  const handleFileSelect = (event: ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files;
    if (files && files.length > 0) {
      uploadFile(files[0]);
    }
  };

  const handleDrop = (event: React.DragEvent<HTMLDivElement>) => {
    event.preventDefault();
    setDragOver(false);

    const files = event.dataTransfer.files;
    if (files && files.length > 0) {
      uploadFile(files[0]);
    }
  };

  const handleDragOver = (event: React.DragEvent<HTMLDivElement>) => {
    event.preventDefault();
    setDragOver(true);
  };

  const handleDragLeave = () => {
    setDragOver(false);
  };

  const handleClick = () => {
    fileInputRef.current?.click();
  };

  return (
    <div className={className}>
      <div
        onClick={handleClick}
        onDrop={handleDrop}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        className={`
          relative border-2 border-dashed rounded-lg p-6
          cursor-pointer transition-colors
          ${dragOver ? "border-blue-500 bg-blue-50 dark:bg-blue-900/20" : "border-gray-300 dark:border-gray-600"}
          ${uploading ? "opacity-50 cursor-not-allowed" : "hover:border-gray-400 dark:hover:border-gray-500"}
        `}
      >
        <input
          ref={fileInputRef}
          type="file"
          accept={accept}
          onChange={handleFileSelect}
          disabled={uploading}
          className="hidden"
        />
        
        <div className="flex flex-col items-center text-center">
          <svg
            className="w-12 h-12 text-gray-400 dark:text-gray-500 mb-3"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
            />
          </svg>
          
          {uploading ? (
            <p className="text-sm text-gray-600 dark:text-gray-400">
              Uploading...
            </p>
          ) : (
            <>
              <p className="text-sm text-gray-600 dark:text-gray-400 mb-1">
                <span className="font-semibold text-blue-600 dark:text-blue-400">
                  Click to upload
                </span>
                {" or drag and drop"}
              </p>
              <p className="text-xs text-gray-500 dark:text-gray-500">
                Max size: {Math.round(maxSize / 1024 / 1024)}MB
              </p>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
