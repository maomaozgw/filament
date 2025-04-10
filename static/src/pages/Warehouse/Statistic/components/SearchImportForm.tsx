import { DeleteRowOutlined, DoubleLeftOutlined, FileAddOutlined, InsertRowBelowOutlined, LeftOutlined, RightOutlined } from "@ant-design/icons";
import { ProForm, ProFormGroup, ProFormSelect } from "@ant-design/pro-components";
import { Button, Form, Row, Space } from "antd";
import { useState } from "react";
import ImportForm from "../../Filament/components/ImportForm";

interface Props {
    onFinished: (changed: boolean) => void;
}

const SearchForm = (props: Props) => {
    const [stockInVisiable, setStockInVisiable] = useState<boolean>(false);
    const [batchStockInVisiable, setBatchStockInVisiable] = useState<boolean>(false);
    const { onFinished } = props;
    return (
        <>
            <ProForm
                name="search_form" layout="inline"
                submitter={{
                    render: () => (
                        <Space>
                            <Button type="primary" htmlType="submit" icon={<LeftOutlined />} onClick={() => setStockInVisiable(true)}>
                                入库
                            </Button>
                            {/* <Button type="primary" htmlType="submit" icon={<RightOutlined />}>
                            出库
                        </Button> */}
                            <Button type="primary" htmlType="submit" icon={<DoubleLeftOutlined />} onClick={() => setBatchStockInVisiable(true)}>
                                批量入库
                            </Button>
                        </Space >

                    ),
                }}
            >
                <ProFormGroup>
                    {/* <ProFormSelect label="品牌" name="brand">
                </ProFormSelect>
                <ProFormSelect label="颜色" name="color">
                </ProFormSelect>
                <ProFormSelect label="类型" name="type">
                </ProFormSelect> */}
                </ProFormGroup>
            </ProForm >
            <ImportForm
                modalVisible={batchStockInVisiable}
                onFinished={(changed) => {
                    setBatchStockInVisiable(false)
                    onFinished(changed)
                }}
            />
        </>
    )
}

export default SearchForm;