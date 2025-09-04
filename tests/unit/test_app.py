# 单元测试 - 应用测试
import unittest
import sys
import os

# 添加项目根目录到路径
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '../../backend'))

class TestApp(unittest.TestCase):
    """应用基础测试"""
    
    def setUp(self):
        """测试前准备"""
        pass
    
    def test_app_import(self):
        """测试应用导入"""
        try:
            from app import app
            self.assertIsNotNone(app)
        except ImportError as e:
            self.fail(f"无法导入应用: {e}")
    
    def test_config_loading(self):
        """测试配置加载"""
        try:
            from config import get_config
            config = get_config()
            self.assertIsNotNone(config)
        except ImportError as e:
            self.fail(f"无法导入配置: {e}")

if __name__ == '__main__':
    unittest.main()
