import { ModelFilament } from "@/models/api";
import { UploadOutlined } from "@ant-design/icons";
import { EditableProTable, ProForm, ProFormGroup, StepsForm } from "@ant-design/pro-components";
import { Button, Modal, Upload } from "antd";
import TextArea from "antd/es/input/TextArea";
import { useState } from "react";
import { filamentColumns } from "./Define";
import { createImport } from "@/services/importer/Import";

interface ImportProps {
    modalVisible: boolean;
    onFinished: (changed: boolean) => void;
}

const ImportForm: React.FC<ImportProps> = (props) => {
    const [rawVal, setRawVal] = useState<string>("");
    const [importData, setImportData] = useState<readonly ModelFilament[]>([]);
    const { modalVisible, onFinished } = props;
    return (
        <Modal
            title="导入"
            open={modalVisible}
            onOk={async () => {
                const result = await createImport("filament", importData);
                onFinished(result)
            }}
            onCancel={() => onFinished(false)}
            maskClosable={false}
            destroyOnClose={true}
            width={900}
        >
            <EditableProTable<ModelFilament>
                key="preview-table"
                rowKey="id"
                loading={false}
                value={importData}
                toolBarRender={() => {
                    return [
                        <Upload
                            key="upload-csv"
                            name="upload-filaments-csv"
                            accept={".csv"}
                            beforeUpload={async (info) => {
                                if (!info) {
                                    return false;
                                }
                                const reader = new FileReader();
                                reader.readAsText(info);
                                reader.onload = () => {
                                    const data: ModelFilament[] = [];
                                    reader.result?.toString().split("\n").forEach((line, idx) => {
                                        const [brand, type, color, price, dt, quantity] = line.split(",");
                                        if (isNaN(parseFloat(price))) {
                                            return;
                                        }
                                        const datetime = new Date(dt);
                                        data.push({
                                            id: idx,
                                            brand: {
                                                name: brand,
                                            },
                                            type: {
                                                name: type,
                                            },
                                            color: {
                                                name: color,
                                            },
                                            price: Math.floor(parseFloat(price) * 100),
                                            created_at: datetime.toISOString(),
                                            quantity: parseInt(quantity),
                                        })
                                    })
                                    data.concat(...importData);
                                    setImportData(data);
                                    return;
                                };
                                return false;

                            }}
                        >
                            <Button icon={<UploadOutlined />}>上传 CSV 文件</Button>
                        </Upload>
                    ]
                }}
                columns={
                    filamentColumns([{
                        title: '操作',
                        valueType: 'option',
                        width: 200,
                        render: (text, record, _, action) => [
                            <a
                                key="editable"
                                onClick={() => {
                                    action?.startEditable?.(record?.id ?? 0, record);
                                }}
                            >
                                编辑
                            </a>,
                            <a
                                key="delete"
                                onClick={() => {
                                }}
                            >
                                删除
                            </a>,
                        ],
                    },])
                }
                recordCreatorProps={{
                    // 每次新增的时候需要Key
                    record: () => ({ id: importData.length + 1 }),
                }}
                onChange={setImportData}
                request={async () => {
                    return {
                        data: [],
                        success: true,
                        total: 0,
                    };
                }}

            />
        </Modal>
    )
}

export default ImportForm;