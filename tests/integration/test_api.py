# 集成测试 - API测试
import unittest
import requests
import time

class TestAPI(unittest.TestCase):
    """API集成测试"""
    
    BASE_URL = "http://localhost:9000"
    
    @classmethod
    def setUpClass(cls):
        """测试类准备"""
        # 等待服务启动
        time.sleep(2)
    
    def test_status_endpoint(self):
        """测试状态端点"""
        try:
            response = requests.get(f"{self.BASE_URL}/api/status", timeout=5)
            self.assertEqual(response.status_code, 200)
            
            data = response.json()
            self.assertIn('status', data)
            self.assertEqual(data['status'], 'running')
            
        except requests.exceptions.ConnectionError:
            self.skipTest("服务器未启动，跳过API测试")
    
    def test_indices_endpoint(self):
        """测试指数列表端点"""
        try:
            response = requests.get(f"{self.BASE_URL}/api/indices", timeout=5)
            self.assertEqual(response.status_code, 200)
            
            data = response.json()
            self.assertIsInstance(data, dict)
            
        except requests.exceptions.ConnectionError:
            self.skipTest("服务器未启动，跳过API测试")
    
    def test_prediction_endpoint(self):
        """测试预测端点"""
        try:
            response = requests.get(f"{self.BASE_URL}/api/predict/sh000001", timeout=10)
            self.assertEqual(response.status_code, 200)
            
            data = response.json()
            self.assertIn('current', data)
            self.assertIn('predicted', data)
            self.assertIn('confidence', data)
            
        except requests.exceptions.ConnectionError:
            self.skipTest("服务器未启动，跳过API测试")

if __name__ == '__main__':
    unittest.main()
