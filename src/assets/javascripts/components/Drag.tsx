import { useCallback } from "preact/hooks";

interface DragProps {
  width: number;
  onResize: (width: number) => void;
}

export default function Drag({ width, onResize }: DragProps) {
  const onMouseDown = useCallback(
    (e: MouseEvent) => {
      const startX = e.clientX;
      const initW = width;

      const onMouseMove = (ev: MouseEvent) => {
        const offset = ev.clientX - startX;
        onResize(initW + offset);
      };

      const onMouseUp = () => {
        document.removeEventListener("mousemove", onMouseMove);
        document.removeEventListener("mouseup", onMouseUp);
      };

      document.addEventListener("mousemove", onMouseMove);
      document.addEventListener("mouseup", onMouseUp);
    },
    [width, onResize],
  );

  return <div className="drag" onMouseDown={onMouseDown} />;
}
