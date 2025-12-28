/**
 * useScreenCapture hook for capturing screen using browser API
 */

export function useScreenCapture() {
  /**
   * Captures the screen and returns a Blob containing the image
   */
  const captureScreen = async (): Promise<Blob> => {
    try {
      // Request screen capture permission
      const stream = await navigator.mediaDevices.getDisplayMedia({
        video: {
          // @ts-ignore - mediaSource is a valid property but not in standard types
          mediaSource: "screen",
        } as MediaTrackConstraints,
      });

      // Create video element to capture the stream
      const video = document.createElement("video");
      video.srcObject = stream;
      video.autoplay = true;

      // Wait for video to be ready (avoid race condition)
      await new Promise<void>((resolve) => {
        if (video.readyState >= 1) { // HAVE_METADATA
          resolve();
        } else {
          video.addEventListener("loadedmetadata", () => {
            resolve();
          }, { once: true });
        }
      });

      // Create canvas to capture frame
      const canvas = document.createElement("canvas");
      canvas.width = video.videoWidth;
      canvas.height = video.videoHeight;

      const ctx = canvas.getContext("2d");
      if (!ctx) {
        throw new Error("Failed to get canvas context");
      }

      // Draw video frame to canvas
      ctx.drawImage(video, 0, 0, canvas.width, canvas.height);

      // Stop all tracks
      stream.getTracks().forEach((track) => track.stop());

      // Convert canvas to blob
      return new Promise<Blob>((resolve, reject) => {
        canvas.toBlob((blob) => {
          if (blob) {
            resolve(blob);
          } else {
            reject(new Error("Failed to create blob from canvas"));
          }
        }, "image/png");
      });
    } catch (error) {
      if (error instanceof Error) {
        if (error.name === "NotAllowedError") {
          throw new Error("Screen capture permission denied");
        } else if (error.name === "NotFoundError") {
          throw new Error("No screen capture source available");
        }
      }
      throw error;
    }
  };

  return {
    captureScreen,
  };
}
