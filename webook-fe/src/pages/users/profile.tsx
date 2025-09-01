import { ProDescriptions } from '@ant-design/pro-components';
import React, { useState, useEffect } from 'react';
import { Button, Spin, message } from 'antd';
import axios from "@/axios/axios";

function Page() {
    const [data, setData] = useState<Profile | null>(null);
    const [isLoading, setLoading] = useState(true);

    useEffect(() => {
        axios.get('/users/profile')
            .then((res) => setData(res.data))
            .catch(() => {
                message.error("获取用户信息失败");
            })
            .finally(() => setLoading(false));
    }, []);

    if (isLoading) {
        return (
            <div style={{ textAlign: "center", paddingTop: 100 }}>
                <Spin size="large" tip="加载中..." />
            </div>
        );
    }

    if (!data) return <p>暂无用户信息</p>;

    return (
        <ProDescriptions<Profile>
            column={1}
            title="个人信息"
            bordered
        >
            <ProDescriptions.Item label="昵称">
                {data.Nickname}
            </ProDescriptions.Item>
            <ProDescriptions.Item label="邮箱">
                {data.Email}
            </ProDescriptions.Item>
            <ProDescriptions.Item label="手机">
                {data.Phone}
            </ProDescriptions.Item>
            <ProDescriptions.Item label="生日" valueType="date">
                {data.Birthday}
            </ProDescriptions.Item>
            <ProDescriptions.Item label="关于我">
                {data.AboutMe}
            </ProDescriptions.Item>
            <ProDescriptions.Item label="操作">
                <div style={{ textAlign: "right", width: "100%" }}>
                    <Button href="/users/edit" type="primary">修改</Button>
                </div>
            </ProDescriptions.Item>
        </ProDescriptions>
    );
}

export default Page;
