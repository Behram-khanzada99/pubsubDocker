# main.tf

provider "kubernetes" {
  config_path = "~/.kube/config"  # Path to your Kubernetes config file
}

# Deploy the 'mydata' PersistentVolumeClaim
resource "kubernetes_persistent_volume_claim" "mydata-t" {
  metadata {
    name = "mydata-t"
    labels = {
      app = "my-go-app-tf"
    }
  }

  spec {
    access_modes = ["ReadWriteOnce"]

    resources {
      requests = {
        storage = "100Mi"
      }
    }
  }
}

# Deploy the 'redis' Deployment
resource "kubernetes_deployment" "redis-tf" {
  metadata {
    name = "redis-tf"
    labels = {
      app = "redis-tf"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "redis-tf"
      }
    }

    template {
      metadata {
        labels = {
          app = "redis-tf"
        }
      }

      spec {
        container {
          name  = "redis-tf"
          image = "redis:latest"

        
        }
      }
    }
  }
}

# Deploy the 'redis' Service
resource "kubernetes_service" "redis-tf-s" {
  metadata {
    name = "redis-tf-s"
    labels = {
      app = "redis-tf"
    }
  }

  spec {
    selector = {
      app = "redis-tf"
    }

    port {
      name       = "6379"
      port       = 6379
      target_port = 6379
    }
  }
}


# Deploy the 'app' Deployment
resource "kubernetes_deployment" "app-tf" {
    timeouts {
    create = "15m"  # Increase this value as needed
  }
  depends_on = [kubernetes_persistent_volume_claim.mydata-t]
  metadata {
    name = "app-tf"
    labels = {
      app = "my-go-app-tf"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "my-go-app-tf"
      }
    }

    template {
      metadata {
        labels = {
          app = "my-go-app-tf"
        }
      }

      spec {
        container {
          name  = "my-go-app-tf"
          image = "docker.io/behramkhanzada/my-go-app:v1"

          
          volume_mount {
            name       = "mydata-t"
            mount_path = "/app/data"
          }
        }

        volume {
          name = "mydata-t"

          persistent_volume_claim {
            claim_name = kubernetes_persistent_volume_claim.mydata-t.metadata[0].name
          }
        }
      }
    }
  }
}


# Deploy the 'app' Service
resource "kubernetes_service" "app-tf-s" {
  metadata {
    name = "app-tf-s"
    labels = {
      app = "my-go-app-tf"
    }
  }

  spec {
    selector = {
      app = "my-go-app-tf"
    }

    port {
      name       = "8081"
      port       = 8081
      target_port = 8080
    }
  }
}

