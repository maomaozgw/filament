import { ModelBrand } from "@/models/api";
import { ActionType, ModalForm, ProForm, ProFormColorPicker, ProFormText } from "@ant-design/pro-components";
import { request } from "@umijs/max";
import { Form } from "antd";
import { AggregationColor } from "antd/es/color-picker/color";

type EditFormProps = {
    currentItem: ModelBrand | undefined,
    actionRef: React.MutableRefObject<ActionType | undefined>,
    updateModalVisible: boolean,
    handleUpdateModalVisible: (visible: boolean) => void,
}

async function create(value: ModelBrand): Promise<boolean> {
    try {
        const { data, code, message } = await request(`/api/v1/meta-data/brands`, {
            method: "POST",
            data: value,
        });
        return true;
    } catch (error) {
        return false;
    }
}

async function update(value: ModelBrand): Promise<boolean> {
    try {
        const { data, code, message } = await request(`/api/v1/meta-data/brands/${value.id}`, {
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
    let title = "编辑品牌";
    let saveFunc = update;
    if (currentItem === undefined) {
        title = "新建拼盘";
        saveFunc = create
        return null;
    }
    
    const [updateForm] = Form.useForm<ModelBrand>();
    return (
        <>
            <ModalForm<ModelBrand>
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
                </ProForm.Group>
            </ModalForm>
        </>
    )
}

export default EditForm;