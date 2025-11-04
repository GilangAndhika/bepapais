package streaming

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/GilangAndhika/bepapais/internal/models"
)

const FFMPEG_PATH = "./ffmpeg.exe" // Pastikan ffmpeg ada di PATH

// Manager mengelola semua proses FFmpeg yang aktif
type Manager struct {
	// Mutex untuk melindungi map, karena akan diakses dari goroutine
	mu      sync.RWMutex
	streams map[string]*exec.Cmd // map[camera_id] -> proses ffmpeg
}

// NewManager membuat instance streaming manager
func NewManager() *Manager {
	return &Manager{
		streams: make(map[string]*exec.Cmd),
	}
}

// StartStream memulai stream FFmpeg untuk satu kamera
func (m *Manager) StartStream(cam models.Camera) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Jika stream sudah berjalan, hentikan dulu
	if cmd, exists := m.streams[cam.ID]; exists {
		log.Printf("[Stream] Menghentikan stream lama untuk: %s", cam.ID)
		if err := cmd.Process.Kill(); err != nil {
			log.Printf("[Stream] Gagal menghentikan proses lama: %v", err)
		}
	}

	// Buat URL RTSP
	rtspURL := fmt.Sprintf("rtsp://%s:%s@%s:%d%s",
		cam.Source.Username,
		cam.Source.Password,
		cam.Source.IP,
		cam.Source.Port,
		cam.Source.Path,
	)

	// Folder output: ./media/[camera_id]
	outputDir := filepath.Join("media", cam.ID)
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Printf("[Stream] Gagal membuat direktori media: %v", err)
		return
	}

	playlistPath := filepath.Join(outputDir, "index.m3u8")

	log.Printf("[Stream] Memulai stream untuk: %s (%s)", cam.Name, cam.ID)

	cmd := exec.Command(
		FFMPEG_PATH,
		"-rtsp_transport", "tcp",
		"-i", rtspURL,
		"-c:v", "copy",
		"-c:a", "copy",
		"-f", "hls",
		"-hls_time", "2",
		"-hls_list_size", "3",
		"-hls_flags", "delete_segments",
		playlistPath,
	)

	cmd.Stderr = os.Stderr // Tampilkan error FFmpeg di console Go

	// Jalankan perintah dalam goroutine baru
	go func() {
		err := cmd.Run()
		if err != nil {
			// Jika error (misal: CVR mati), log error
			log.Printf("[Stream] Error FFmpeg [ %s ]: %v", cam.Name, err)
		}
		
		// Hapus dari map jika proses berhenti
		m.mu.Lock()
		delete(m.streams, cam.ID)
		m.mu.Unlock()
		log.Printf("[Stream] Stream untuk %s telah berhenti.", cam.Name)
	}()

	// Simpan proses di map
	m.streams[cam.ID] = cmd
}

// StopStream menghentikan stream FFmpeg berdasarkan ID kamera
func (m *Manager) StopStream(cameraID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	cmd, exists := m.streams[cameraID]
	if !exists {
		log.Printf("[Stream] Tidak ada stream aktif untuk dihentikan: %s", cameraID)
		return
	}

	log.Printf("[Stream] Menghentikan stream untuk: %s", cameraID)
	if err := cmd.Process.Kill(); err != nil {
		log.Printf("[Stream] Gagal menghentikan proses: %v", err)
	}
	
	delete(m.streams, cameraID)
}