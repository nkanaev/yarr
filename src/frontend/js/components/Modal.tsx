import { useEffect, useRef } from "preact/hooks";
import type { ComponentChildren } from "preact";

interface ModalProps {
  open: boolean;
  onHide: () => void;
  children: ComponentChildren;
}

export default function Modal({ open, onHide, children }: ModalProps) {
  const dialogRef = useRef<HTMLDialogElement>(null);

  useEffect(() => {
    const el = dialogRef.current;
    if (!el) return;
    if (open) {
      if (!el.open) el.showModal();
    } else {
      if (el.open) el.close();
    }
  }, [open]);

  const onBackdropClick = (e: MouseEvent) => {
    if (e.target === dialogRef.current) {
      onHide();
    }
  };

  return (
    <dialog className="c-dialog" ref={dialogRef} onClose={onHide} onClick={onBackdropClick}>
      {children}
    </dialog>
  );
}
