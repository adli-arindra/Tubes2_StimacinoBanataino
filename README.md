# README - Stimacinobanataino Project

## i. Penjelasan Singkat Algoritma DFS dan BFS yang Diimplementasikan

### **DFS (Depth-First Search)**  
DFS adalah algoritma pencarian pada graf yang menjelajahi sejauh mungkin sepanjang cabang sebelum kembali. Algoritma ini bekerja dengan memulai dari sebuah node dan mencoba untuk menyelami setiap cabang hingga mencapai node akhir atau tidak ada lagi node yang bisa dijelajahi. Jika jalur tersebut sudah selesai, algoritma akan kembali ke node sebelumnya dan mencari cabang lain.

### **BFS (Breadth-First Search)**  
BFS adalah algoritma pencarian pada graf yang menjelajahi graf secara level per level. Dimulai dari node awal, BFS akan mengeksplorasi node yang berdekatan terlebih dahulu sebelum melanjutkan ke node yang lebih jauh. Algoritma ini menggunakan antrian untuk melacak node yang perlu dieksplorasi.

## ii. Requirement Program dan Instalasi Tertentu

Untuk menjalankan proyek ini, Anda memerlukan beberapa alat dan dependensi:

### **Requirement:**
- **Docker**: Untuk menjalankan aplikasi dalam container.
- **Node.js**: Versi 18 atau lebih baru (untuk menjalankan frontend).

### **Instalasi:**
1. Pastikan Docker sudah terinstal pada sistem Anda.
2. Pastikan Anda memiliki akses ke repositori backend dan frontend.

## iii. Command atau Langkah-Langkah dalam Meng-compile atau Build Program

### **Clone Repository**
```bash
git clone https://github.com/adli-arindra/Tubes2_StimacinoBanataino.git
```

### **Cara Menjalankan Backend**
```bash
cd backend
docker build -t stimacinobanataino-be
docker run -p 8080:8080 stimacinobanataino-be
```

### **Cara Menjalankan Frontend**
```bash
cd frontend/app
echo "NEXT_PUBLIC_ENDPOINT=http://localhost:8080/search" > .env
docker build -t stimacinobanataino-fe
docker run -p 3000:3000 stimacinobanataino-be
```

## iv. Author
- Muhammad Adli Arindra (1822089)
- Andri Nurdianto (13523145)
- Hasri Fayadh Muqaffa (13523156)
