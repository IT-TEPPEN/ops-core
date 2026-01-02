// Layout components
export { ThreePaneLayout } from "./Layout/ThreePaneLayout";
export type { ThreePaneLayoutProps } from "./Layout/ThreePaneLayout";

export { UI_Form_Field } from "../../../ui/form/FormField";
export type { FormFieldProps } from "../../../ui/form/FormField";

export { VariableForm } from "./Form/VariableForm";
export type { VariableFormProps } from "./Form/VariableForm";

export { ExecutionStepPanel } from "./Form/ExecutionStepPanel";
export type { ExecutionStepPanelProps } from "./Form/ExecutionStepPanel";

export { ScreenCaptureButton } from "./Form/ScreenCaptureButton";
export type { ScreenCaptureButtonProps } from "./Form/ScreenCaptureButton";

export { GroupForm } from "./Form/GroupForm";

export { GroupMemberSelector } from "./Form/GroupMemberSelector";

// Display components
export { AttachmentList } from "./Display/AttachmentList";
export type { AttachmentListProps } from "./Display/AttachmentList";

export { AttachmentUploader } from "./Display/AttachmentUploader";
export type {
  AttachmentUploaderProps,
  AttachmentResponse,
} from "./Display/AttachmentUploader";

export { AttachmentViewer } from "./Display/AttachmentViewer";
export type { AttachmentViewerProps } from "./Display/AttachmentViewer";

// Display components - UI components moved to src/ui
export { Card } from "@/ui/Card";
export type { CardProps } from "@/ui/Card";

export { Tag } from "@/ui/Tag";
export type { TagProps } from "@/ui/Tag";

export { ViewHistoryList } from "./Display/ViewHistoryList";
export type { ViewHistoryListProps } from "./Display/ViewHistoryList";

export { DocumentStatistics } from "./Display/DocumentStatistics";
export type { DocumentStatisticsProps } from "./Display/DocumentStatistics";

export { PopularDocumentsList } from "./Display/PopularDocumentsList";
export type { PopularDocumentsListProps } from "./Display/PopularDocumentsList";

export { RecentViewsList } from "./Display/RecentViewsList";
export type { RecentViewsListProps } from "./Display/RecentViewsList";

export { StatisticsChart } from "./Display/StatisticsChart";
export type { StatisticsChartProps } from "./Display/StatisticsChart";

export { ExecutionRecordList } from "./Display/ExecutionRecordList";
export type { ExecutionRecordListProps } from "./Display/ExecutionRecordList";

export { GroupList } from "./Display/GroupList";

export { GroupMemberList } from "./Display/GroupMemberList";

export * from "./ProtectedRoute";
