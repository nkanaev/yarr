import type { Folder } from "../api-types";
import { useI18n } from "../i18n";
import Dropdown from "./Dropdown";
import Icon from "./Icon";

interface DropdownFolderProps {
  folder: Folder;
  onRename: (folder: Folder) => void;
  onDelete: (folder: Folder) => void;
}

export default function DropdownFolder({ folder, onRename, onDelete }: DropdownFolderProps) {
  const { t } = useI18n();

  return (
    <Dropdown
      toggleClass="c-button-link px-2"
      drop="right"
      title={t("folder_settings")}
      button={<Icon name="more-horizontal" />}>
      <div className="c-dropdown-header text-break line-clamp-1" role="heading" aria-level={2}>
        {folder.title}
      </div>
      <button
        type="button"
        className="c-dropdown-item w-100 text-start d-flex gap-1"
        onClick={() => onRename(folder)}>
        <Icon className="me-1" name="edit" />
        {t("rename")}
      </button>
      <div className="c-dropdown-divider" />
      <button
        type="button"
        className="c-dropdown-item w-100 text-start text-danger d-flex gap-1"
        onClick={() => onDelete(folder)}>
        <Icon className="me-1" name="trash" />
        {t("delete")}
      </button>
    </Dropdown>
  );
}
