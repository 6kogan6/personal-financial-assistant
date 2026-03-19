
export default function ConfirmDialog({
  isOpen,
  title,
  message,
  confirmText = "Удалить",
  cancelText = "Отмена",
  onConfirm,
  onCancel,
  isLoading = false,
}) {
  if (!isOpen) return null;

  return (
    <div style={overlayStyle}>
      <div style={modalStyle}>
        <h3 style={titleStyle}>{title}</h3>
        <p style={messageStyle}>{message}</p>

        <div style={actionsStyle}>
          <button
            type="button"
            onClick={onCancel}
            disabled={isLoading}
            style={cancelButtonStyle}
          >
            {cancelText}
          </button>

          <button
            type="button"
            onClick={onConfirm}
            disabled={isLoading}
            style={confirmButtonStyle}
          >
            {isLoading ? "Удаление..." : confirmText}
          </button>
        </div>
      </div>
    </div>
  );
}

const overlayStyle = {
  position: "fixed",
  inset: 0,
  backgroundColor: "rgba(0, 0, 0, 0.45)",
  display: "flex",
  alignItems: "center",
  justifyContent: "center",
  zIndex: 9999,
  padding: "16px",
};

const modalStyle = {
  width: "100%",
  maxWidth: "420px",
  backgroundColor: "#ffffff",
  borderRadius: "16px",
  padding: "24px",
  boxShadow: "0 10px 30px rgba(0, 0, 0, 0.2)",
};

const titleStyle = {
  margin: 0,
  marginBottom: "12px",
  fontSize: "22px",
  fontWeight: 700,
  color: "#111827",
};

const messageStyle = {
  margin: 0,
  marginBottom: "20px",
  fontSize: "16px",
  lineHeight: 1.5,
  color: "#374151",
};

const actionsStyle = {
  display: "flex",
  justifyContent: "flex-end",
  gap: "12px",
};

const cancelButtonStyle = {
  padding: "10px 16px",
  borderRadius: "10px",
  border: "1px solid #d1d5db",
  backgroundColor: "#ffffff",
  color: "#111827",
  cursor: "pointer",
  fontSize: "14px",
  fontWeight: 600,
};

const confirmButtonStyle = {
  padding: "10px 16px",
  borderRadius: "10px",
  border: "none",
  backgroundColor: "#dc2626",
  color: "#ffffff",
  cursor: "pointer",
  fontSize: "14px",
  fontWeight: 600,
};
