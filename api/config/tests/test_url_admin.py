from django.test import TestCase, Client

class TestUrlAdmin(TestCase):
    """ Test admin view ("/admin") """
    
    def setUp(self):
        self.client = Client()
        
    def test_admin_view_status_code(self):
        """ Test admin view retourne status code 200 """
        response = self.client.get('/admin/', follow=True)
        self.assertEqual(response.status_code, 200)
        
    