import { ModelFilament } from "@/models/api";
import { stockIn } from "@/services/warehouse/Filament";
import { ModalForm, ProForm, ProFormDigit, ProFormMoney, ProFormText } from "@ant-design/pro-components";
import React, { PropsWithChildren } from "react";

interface CreateProps {
    modalVisible: boolean;
    onFinished: (changed: boolean) => void;
}


const CreatForm: React.FC<CreateProps> = (props) => {
    const { modalVisible, onFinished } = props;
    return (
        <ModalForm<ModelFilament>
            title="入库耗材"
            autoFocusFirstInput
            open={modalVisible}
            modalProps={{
                destroyOnClose: true,
                width: 420,
                maskClosable: false,
                onCancel: () => onFinished(false),
            }}
            onFinish={async (value) => {
                const result = await stockIn(value);
                onFinished(result);
                return result;
            }}
        >   <ProForm.Group>
                <ProFormText width="md" name={["type", "name"]} label="类型" />
                <ProFormText width="md" name={["brand", "name"]} label="品牌" />
                <ProFormText width="md" name={["color", "name"]} label="颜色" />
            </ProForm.Group>
            <ProForm.Group>
                <ProFormDigit width="sm" name="quantity" label="数量" min={0} />
                <ProFormMoney width="sm" name="price" label="价格" min={0} />

            </ProForm.Group>
        </ModalForm>
    );
}

export default CreatForm;