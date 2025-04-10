import { request } from "@umijs/max";
import { ModelColor } from "@/models/api"
import {
    ActionType, ModalForm, PageContainer, ProDescriptionsItemProps, ProForm,
    ProFormDigit, ProFormText, ProTable, RequestData, ProFormColorPicker
} from "@ant-design/pro-components"
import { useRef, useState } from "react";
import { Button, Form, Modal } from "antd";
import { AggregationColor } from "antd/es/color-picker/color";
import EditForm from "./Components/EditForm";

async function queryColors(params): Promise<Partial<RequestData<ModelColor>>> {
    const { current, pageSize, ...filters } = params;
    const { data, pager } = await request("/api/v1/meta-data/colors", {
        method: "GET",
        params: {
            page: current,
            size: pageSize,
            ...filters,
        },
    });
    return {
        data,
        success: true,
        total: pager.total,
    }
}


const Color: React.FC = () => {
    const actionRef = useRef<ActionType>();
    const [updateModalVisible, handleUpdateModalVisible] = useState<boolean>(false);
    const [currentItem, setCurrentItem] = useState<ModelColor>({});
    const columns: ProDescriptionsItemProps<ModelColor>[] = [
        {
            title: "id",
            dataIndex: "id",
            valueType: "digit",
            hideInTable: true,
            hideInForm: true,
            hideInSearch: true,
        },
        {
            title: "名称",
            dataIndex: "name",
            valueType: "text",
        },
        {
            title: "颜色",
            dataIndex: "rgba",
            valueType: "color",
            hideInSearch: true,
        },
        {
            title: "更新时间",
            dataIndex: "updated_at",
            valueType: "dateTime",
            hideInForm: true,
            hideInSearch: true,
        },
        {
            title: "操作",
            dataIndex: "option",
            valueType: "option",
            render: (_, record) => (
                <>
                    <a onClick={() => {
                        handleUpdateModalVisible(true);
                        setCurrentItem(record);
                    }}>编辑</a>
                </>
            )
        }
    ]
    return (
        <PageContainer
            header={{
                title: "颜色管理",
            }}
        >
            <ProTable<ModelColor>
                rowKey="id"
                actionRef={actionRef}
                request={queryColors}
                columns={columns}
                toolBarRender={() => [
                    <Button
                        key="1"
                        type="primary"
                        onClick={() => {
                            setCurrentItem({});
                            handleUpdateModalVisible(true)
                        }}
                    >
                        新建
                    </Button>,
                ]}
            />
            <EditForm
                currentItem={currentItem}
                actionRef={actionRef}
                updateModalVisible={updateModalVisible}
                handleUpdateModalVisible={handleUpdateModalVisible}
            ></EditForm>
        </PageContainer>
    )
}

export default Color;