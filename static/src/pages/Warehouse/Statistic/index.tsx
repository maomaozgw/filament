import { ModelFilament, ModelStatistic, ModelStatisticValue } from '@/models/api';
import { Datum, Pie, Sunburst } from '@ant-design/charts';
import { DatabaseOutlined, DownOutlined, MoneyCollectOutlined, UpOutlined } from '@ant-design/icons';
import { PageContainer, StatisticCard, StatisticProps } from '@ant-design/pro-components';
import { request } from '@umijs/max';
import { Divider, Form, message, Space } from 'antd';
import { getStatusClassNames } from 'antd/es/_util/statusUtils';
import { useEffect, useState } from 'react';
import SearchForm from './components/SearchImportForm';

interface pieItem {
    type: string;
    value: number;
    precent?: number;
}

const getSatisticSummaryProps = (item: ModelStatisticValue): StatisticProps => {
    switch (item.name) {
        case "Current":
            return {
                icon: <DatabaseOutlined />,
                value: item.value,
                title: "当前库存",
                suffix: "卷",
                style: {
                    color: "#1890ff",
                }
            }
        case "Stock In":
            return {
                icon: <DownOutlined />,
                value: item.value,
                title: "一共购买",
                suffix: "卷",
                style: {
                    color: "#00ff00",
                }
            }
        case "Stock Out":
            return {
                icon: <UpOutlined />,
                value: item.value,
                title: "消耗",
                suffix: "卷",
                style: {
                    color: "#ff0000",
                }
            }
        case "Total Cost":
            return {
                icon: <MoneyCollectOutlined />,
                value: (item.value ?? 0) / 100,
                prefix: '¥',
                title: "总共花费",
                style: {
                    color: "#ff0000",
                }
            }
        default:
            return {
                value: item.value,
                title: item.name,
            }
    }
}

const getPieStatisticCard = (item: ModelStatistic, title: string) => {
    if (item == undefined || (item.values ?? []).length == 0) {
        return <></>
    }
    return <StatisticCard
        key={title}
        title={title}
        chart={
            <Pie
                data={item.values}
                angleField="value"
                colorField="name"
                innerRadius={0.6}
                legend={false}
                style={{
                    stroke: '#fff',
                    inset: 1,
                    radius: 10,
                }}
                scale={{
                    color: {
                        palette: 'spectral',
                        offset: (t) => t * 0.8 + 0.1,
                    },
                }}
                tooltip={(d) => {
                    const result = {
                        name: `${d.name}`,
                        title: {
                            value: d.name,
                        },
                        value: `数量: ${d.value}`,
                    }
                    return result;
                }}
                label={{
                    position: 'inside',
                    formatter: (t: string, d: Datum) => {
                        return `${d.name}: ${d.value}`;
                    },
                }}
            />
        }
    />;
}

const getSunburstStatisticCard = (item: ModelStatistic, title: string) => {
    if (item == undefined || (item.values ?? []).length == 0) {
        return <></>;
    }
    const data = {
        "name": title,
        "children": item.values,
    }
    return (<StatisticCard
        key={title}
        title={title}
        chart={
            <Sunburst
                valueField="value"
                colorField="name"
                data={{
                    type: "inline",
                    value: data,
                }}
                label={{
                    formatter: (t: string, d: Datum) => {
                        return `${d.name}: ${d.value}`;
                    },
                }}
                legend={false}
                animate={{
                    enter: { type: 'waveIn' }
                }}
                innerRadius={0}
            />
        } />);
}

const queryStatistic = async (params) => {
    const { current, pageSize, ...filters } = params;
    const res = await request("/api/v1/warehouse/statistic", {
        method: "GET",
        params: {
            page: current,
            size: pageSize,
            ...filters,
        },
    });
    return res;
}

function countBy(data: ModelFilament[], totalCnt: number, keyGetter: (item: ModelFilament) => string): pieItem[] {
    return data.reduce((prev, cur, _, total) => {
        const key = keyGetter(cur);
        if (!total) {
            return [];
        }
        const index = prev.findIndex((item) => item.type === key);
        if (index === -1) {
            const item = {
                type: key,
                value: cur.quantity ?? 0,
                precent: (cur.quantity ?? 0) / totalCnt * 100,
            }
            prev.push(item);
            console.trace("pushed", item);
        } else {
            prev[index].value += cur.quantity ?? 0;
            prev[index].precent = (prev[index].value) / totalCnt * 100;
            console.trace(prev[index].type, "precent", prev[index].precent);
        }
        return prev;
    }, [] as pieItem[])
}

const Statistic: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [total, setTotal] = useState(0);
    const [summary, setSummary] = useState<ModelStatistic>({});
    const [byBrand, setByBrand] = useState<ModelStatistic>({});
    const [byColor, setByColor] = useState<ModelStatistic>({});
    const [byType, setByType] = useState<ModelStatistic>({});
    const fetchData = async () => {
        queryStatistic({}).then((res) => {
            if (res.code != 0) {
                message.error(res.message);
                return;
            }
            const data = res.data as ModelStatistic[];
            data.forEach((item) => {
                switch (item.kind) {
                    case "summary":
                        setSummary(item);
                        break;
                    case "Pie":
                        switch (item.title) {
                            case "Brand":
                                setByBrand(item);
                                break;
                            case "Color":
                                setByColor(item);
                                break;
                            case "Type":
                                setByType(item);
                                break;
                            default:
                                console.error("unknown title", item.title);
                        }
                        break;
                    case "Sunburst":
                        setByType(item);
                        break;
                    default:
                        console.error("unknown statistic kind", item.kind);
                        break;
                }
            })
        }
        ).catch((err) => {
            message.error(err.message);
        }).finally(() => {
        });
    }
    useEffect(() => {
        fetchData();
        setLoading(false);
    }, [loading])
    return (
        <PageContainer>
            <SearchForm onFinished={(changed) => changed ? setLoading(true) : null} />
            <Divider />
            <StatisticCard.Group direction='row'>
                {
                    summary?.values && (summary.values ?? []).length > 0 &&
                    summary.values.map((item) => {
                        return (<StatisticCard
                            key={item.name}
                            statistic={getSatisticSummaryProps(item)}
                        />)
                    })
                }

            </StatisticCard.Group>
            <StatisticCard.Group direction='row'>
                {getPieStatisticCard(byBrand, "耗材品牌统计")}
                {getPieStatisticCard(byColor, "耗材颜色统计")}
                {getSunburstStatisticCard(byType, "耗材类型统计")}
            </StatisticCard.Group>
        </PageContainer >
    );
}

export default Statistic;