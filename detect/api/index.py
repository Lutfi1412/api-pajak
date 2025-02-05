from flask import Flask, request, jsonify
import requests
import psycopg2
import os
from dotenv import load_dotenv
import tempfile
from datetime import datetime
import pytz

# Load konfigurasi dari file .env
load_dotenv()

app = Flask(__name__)

# Konfigurasi Database
DB_CONFIG = {
    "user": os.getenv("DB_USER"),
    "password": os.getenv("DB_PASSWORD"),
    "host": os.getenv("DB_HOST"),
    "port": os.getenv("DB_PORT"),
    "dbname": os.getenv("DB_NAME"),
}

# Konfigurasi API Plate Recognizer
API_URL = os.getenv("API_URL")
API_KEY = os.getenv("API_KEY")


# Fungsi untuk mendapatkan keterangan status pajak kendaraan
def get_keterangan(db_conn, plat, tanggal, jenis):
    query_pemilik = """
        SELECT tenggat_thn, tenggat_bln 
        FROM data_pemilik 
        WHERE plat = %s AND jenis = %s
    """
    cursor = db_conn.cursor()
    
    try:
        cursor.execute(query_pemilik, (plat, jenis))
        result = cursor.fetchone()
        
        if result is None:
            return "Plat nomor tidak terdaftar"
        
        tenggat_thn, tenggat_bln = result

        # Konversi tanggal sekarang
        tanggal_obj = datetime.strptime(tanggal, "%Y-%m-%d")
        tahun = int(tanggal_obj.strftime("%y"))  # Ambil 2 digit terakhir tahun
        bulan = int(tanggal_obj.strftime("%m"))  # Ambil bulan

        # Cek status pajak berdasarkan tahun dan bulan
        if tahun > tenggat_thn:
            return "Belum bayar pajak"
        elif tahun == tenggat_thn and bulan > tenggat_bln:
            return "Belum bayar pajak"
        else:
            return "Sudah bayar pajak"

    except psycopg2.Error as e:
        return f"Terjadi kesalahan saat validasi data: {str(e)}"
    
    finally:
        cursor.close()


# Fungsi untuk mendeteksi kendaraan dan plat nomor
def detect_vehicles_and_plates(image_path):
    try:
        with open(image_path, "rb") as image_file:
            response = requests.post(
                API_URL,
                files={"upload": image_file},
                headers={"Authorization": f"Token {API_KEY}"}
            )
        response.raise_for_status()  # Memastikan status code adalah 200
        result = response.json()
    except requests.RequestException as e:
        raise Exception(f"API Request Error: {str(e)}")

    detections = []

    for detection in result.get("results", []):
        vehicle_type = detection.get("vehicle", {}).get("type", "Tidak terdeteksi")
        plate = detection.get("plate", "Tidak terdeteksi")

        # Menyederhanakan jenis kendaraan
        if vehicle_type.lower() in ["sedan", "suv", "pickup", "van", "minivan", "truck", "bus", "ambulance", "taxi"]:
            vehicle_type = "mobil"
        elif vehicle_type.lower() == "motorcycle":
            vehicle_type = "motor"
        else:
            vehicle_type = "Tidak terdeteksi"

        # Hanya tambahkan jika jenis kendaraan dan plat terdeteksi
        if vehicle_type != "Tidak terdeteksi" and plate != "Tidak terdeteksi":
            detections.append({"jenis": vehicle_type, "plat": plate})

    return detections


# Fungsi untuk menyimpan hasil ke database dengan waktu dan keterangan
def save_to_database(detections):
    conn = None
    cursor = None
    try:
        conn = psycopg2.connect(**DB_CONFIG)
        cursor = conn.cursor()

        # Ambil waktu sekarang dengan timezone Asia/Jakarta
        jakarta_timezone = pytz.timezone("Asia/Jakarta")
        now = datetime.now(jakarta_timezone)
        tanggal = now.strftime("%Y-%m-%d")
        jam = now.strftime("%H:%M:%S")

        for detection in detections:
            jenis = detection["jenis"]
            plat = detection["plat"]

            # Tentukan logika keterangan berdasarkan tanggal, plat, dan jenis kendaraan
            keterangan = get_keterangan(conn, plat, tanggal, jenis)

            query = """
                INSERT INTO laporan (tanggal, jam, jenis, plat, keterangan)
                VALUES (%s, %s, %s, %s, %s)
            """
            cursor.execute(query, (tanggal, jam, jenis, plat, keterangan))

        conn.commit()
    except psycopg2.DatabaseError as e:
        raise Exception(f"Database Error: {str(e)}")
    finally:
        if cursor:
            cursor.close()
        if conn:
            conn.close()


# Endpoint untuk menerima file dan melakukan deteksi
@app.route('/detect', methods=['POST'])
def detect():
    if 'image' not in request.files:
        return jsonify({"error": "No image file provided"}), 400

    image_file = request.files['image']

    # Validasi format file
    if not image_file.filename.lower().endswith(('.png', '.jpg', '.jpeg')):
        return jsonify({"error": "Invalid file format. Please upload an image."}), 400

    # Simpan file sementara di direktori yang sesuai
    with tempfile.NamedTemporaryFile(delete=False, suffix=".jpg") as temp_file:
        image_file.save(temp_file.name)
        temp_path = temp_file.name

    try:
        # Jalankan deteksi
        detections = detect_vehicles_and_plates(temp_path)

        # Simpan ke database
        save_to_database(detections)

        # Kembalikan respon
        return jsonify({"status": "success", "data": detections}), 200
    except Exception as e:
        return jsonify({"error": str(e)}), 500
    finally:
        # Hapus file sementara
        if os.path.exists(temp_path):
            os.remove(temp_path)


# Menjalankan aplikasi di lokal
if __name__ == "__main__":
    app.run(debug=True)
