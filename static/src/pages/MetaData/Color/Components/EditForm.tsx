import { ModelColor } from "@/models/api";
import { ActionType, ModalForm, ProForm, ProFormColorPicker, ProFormText } from "@ant-design/pro-components";
import { request } from "@umijs/max";
import { Form } from "antd";
import { AggregationColor } from "antd/es/color-picker/color";

type EditFormProps = {
    currentItem: ModelColor | undefined,
    actionRef: React.MutableRefObject<ActionType | undefined>,
    updateModalVisible: boolean,
    handleUpdateModalVisible: (visible: boolean) => void,
}

async function createColor(value: ModelColor): Promise<boolean> {
    try {
        const { data, code, message } = await request(`/api/v1/meta-data/colors`, {
            method: "POST",
            data: value,
        });
        return true;
    } catch (error) {
        return false;
    }
}

async function updateColor(value: ModelColor): Promise<boolean> {
    try {
        const { data, code, message } = await request(`/api/v1/meta-data/colors/${value.id}`, {
            method: "PUT",
            data: value,
        });
        return true;
    } catch (error) {
        return false;
    }
}

const EditForm: React.FC<EditFormProps> = (props) => {
    const { currentItem, actionRef, updateModalVisible, handleUpdateModalVisible } = props;
    if (!updateModalVisible) {
        return null;
    }
    let title = "编辑颜色";
    let saveFunc = updateColor;
    if (currentItem === undefined) {
        title = "新建颜色";
        saveFunc = createColor
        return null;
    }
    
    const [updateForm] = Form.useForm<ModelColor>();
    return (
        <>
            <ModalForm<ModelColor>
                title={title}
                autoFocusFirstInput
                clearOnDestroy
                form={updateForm}
                open={updateModalVisible}
                onOpenChange={(open) => {
                    handleUpdateModalVisible(open);
                    if (!open) {
                        updateForm.resetFields();
                    }
                }}
                initialValues={
                    {
                        ...currentItem,
                    }
                }
                modalProps={{
                    destroyOnClose: true,
                    width: 420,
                }}
                onFinish={async (value) => {
                    value.id = currentItem?.id || 0;
                    if (value.rgba instanceof AggregationColor) {
                        value.rgba = value.rgba.toHexString()
                    }
                    const result = await saveFunc(value);
                    if (result) {
                        actionRef.current?.reload();
                        updateForm.resetFields();
                    }
                    return result;
                }}
            >
                <ProForm.Group>
                    <ProFormText width="md" name="name" label="名称" />
                    <ProFormColorPicker width="sm" name="rgba" label="颜色" min={0} />
                </ProForm.Group>
            </ModalForm>
        </>
    )
}

export default EditForm;