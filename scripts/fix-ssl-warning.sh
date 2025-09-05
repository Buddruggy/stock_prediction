#!/bin/bash

# 解决Chrome SSL警告的快速脚本

echo "🔧 解决Chrome SSL警告"
echo "======================"
echo ""
echo "您看到的'您的连接不是私密连接'警告是因为使用了自签名证书。"
echo ""
echo "解决方案选择:"
echo "1) 继续使用自签名证书（测试环境）"
echo "2) 配置Let's Encrypt免费证书（生产环境推荐）"
echo "3) 查看详细说明"
echo ""

read -p "请选择 (1-3): " choice

case $choice in
    1)
        echo ""
        echo "✅ 继续使用自签名证书"
        echo ""
        echo "在Chrome中:"
        echo "1. 点击'高级'"
        echo "2. 点击'继续访问 gogotou.cn（不安全）'"
        echo "3. 网站将正常加载"
        echo ""
        echo "⚠️  注意: 每次访问都需要手动确认"
        ;;
    2)
        echo ""
        echo "🔐 配置Let's Encrypt免费证书"
        echo ""
        echo "前提条件:"
        echo "1. 域名 gogotou.cn 已解析到服务器IP"
        echo "2. 服务器80端口已开放"
        echo "3. 服务器443端口已开放"
        echo ""
        read -p "确认满足以上条件？(y/n): " confirm
        if [[ $confirm == "y" || $confirm == "Y" ]]; then
            echo ""
            echo "开始配置Let's Encrypt证书..."
            make setup-letsencrypt
        else
            echo "请先满足前提条件后再配置"
        fi
        ;;
    3)
        echo ""
        echo "📚 详细说明"
        echo ""
        echo "问题原因:"
        echo "- 自签名证书不被浏览器信任"
        echo "- Chrome的安全机制阻止访问"
        echo ""
        echo "解决方案对比:"
        echo ""
        echo "自签名证书:"
        echo "✅ 快速配置"
        echo "❌ 浏览器警告"
        echo "❌ 用户信任度低"
        echo "❌ 每次访问需确认"
        echo ""
        echo "Let's Encrypt证书:"
        echo "✅ 免费"
        echo "✅ 浏览器完全信任"
        echo "✅ 绿色锁图标"
        echo "✅ 用户信任度高"
        echo "❌ 需要域名解析"
        echo "❌ 需要服务器配置"
        echo ""
        echo "推荐: 生产环境使用Let's Encrypt证书"
        ;;
    *)
        echo "❌ 无效选择"
        ;;
esac
