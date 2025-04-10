import { ModelFilament } from "@/models/api";
import { stockIn, stockOut } from "@/services/warehouse/Filament";
import { ModalForm, ProForm, ProFormDigit, ProFormMoney, ProFormText } from "@ant-design/pro-components";
import { Form, Modal } from "antd";
import React, { Dispatch, PropsWithChildren, SetStateAction } from "react";

const stockInKey = "stock-in";

interface StockActionProps {
    modalVisible: boolean;
    filament: ModelFilament;
    stockActionType: "stock-in" | "stock-out";
    onFinish: (changed: Boolean) => void;
}

const ActionForm: React.FC<StockActionProps> = (props) => {
    const { modalVisible, filament, stockActionType, onFinish } = props;
    const stockActionName = stockActionType == stockInKey ? "入库" : "出库";
    return (
        <ModalForm<ModelFilament>
            title={stockActionName}
            autoFocusFirstInput
            open={modalVisible}
            initialValues={
                {
                    ...filament,
                    quantity: 1,
                }
            }
            modalProps={{
                destroyOnClose: true,
                width: 420,
                maskClosable: false,
                onCancel: () => onFinish(false),
            }}
            onFinish={async (value) => {
                value.id = filament.id;
                const result = stockActionType == stockInKey ? await stockIn(value) : await stockOut(value);
                onFinish(result);
                return result;
            }}
        >   <ProForm.Group>
                <ProFormText disabled width="md" name={["type", "name"]} label="类型" />
                <ProFormText disabled width="md" name={["brand", "name"]} label="品牌" />
                <ProFormText disabled width="md" name={["color", "name"]} label="颜色" />
            </ProForm.Group>
            <ProForm.Group>
                <ProFormDigit width="sm" name="quantity" label="数量" min={0} />
                {stockActionType == stockInKey ? <ProFormMoney width="sm" name="price" label="价格" min={0} /> : null}

            </ProForm.Group>
        </ModalForm>
    );
}

export default ActionForm;