package handlers

import (
	"bobastream/internal/services"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type StreamHandler struct {
	videoService *services.VideoService
}

func NewStreamHandler(videoService *services.VideoService) *StreamHandler {
	return &StreamHandler{videoService: videoService}
}

// ShowPlayer shows video player page with thumbnail
func (h *StreamHandler) ShowPlayer(c *fiber.Ctx) error {
	token := c.Params("token")

	// Get video by wrapper token
	video, err := h.videoService.GetVideoByWrapperToken(token)
	if err != nil {
		return c.Status(404).SendString(`
			<!DOCTYPE html>
			<html>
			<head><title>Video Not Found</title></head>
			<body style="background:#000;color:#fff;text-align:center;padding:100px;font-family:Arial;">
				<h1>❌ Video Not Found</h1>
				<p>This video doesn't exist or has been removed.</p>
				<a href="/" style="color:#4a9eff;text-decoration:none;">← Back to Home</a>
			</body>
			</html>
		`)
	}

	// Build player HTML
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s - BOBA STREAM</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            background: #0f0f0f;
            color: #fff;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
        }
        .container {
            max-width: 1400px;
            margin: 0 auto;
            padding: 20px;
        }
        .back-link {
            display: inline-block;
            padding: 10px 20px;
            background: rgba(255,255,255,0.1);
            color: white;
            text-decoration: none;
            border-radius: 8px;
            margin-bottom: 20px;
            transition: background 0.2s;
        }
        .back-link:hover {
            background: rgba(255,255,255,0.2);
        }
        .video-wrapper {
            position: relative;
            background: #000;
            border-radius: 12px;
            overflow: hidden;
            box-shadow: 0 8px 32px rgba(0,0,0,0.6);
        }
        .thumbnail {
            width: 100%%;
            display: block;
            cursor: pointer;
        }
        .play-button {
            position: absolute;
            top: 50%%;
            left: 50%%;
            transform: translate(-50%%, -50%%);
            width: 80px;
            height: 80px;
            background: rgba(255,0,0,0.8);
            border-radius: 50%%;
            display: flex;
            align-items: center;
            justify-content: center;
            cursor: pointer;
            transition: all 0.3s;
        }
        .play-button:hover {
            background: rgba(255,0,0,1);
            transform: translate(-50%%, -50%%) scale(1.1);
        }
        .play-button::after {
            content: '';
            width: 0;
            height: 0;
            border-left: 25px solid white;
            border-top: 15px solid transparent;
            border-bottom: 15px solid transparent;
            margin-left: 5px;
        }
        video {
            width: 100%%;
            display: none;
        }
        video.active {
            display: block;
        }
        .info {
            margin-top: 20px;
        }
        .title {
            font-size: 24px;
            font-weight: 600;
            margin-bottom: 12px;
        }
        .meta {
            color: #aaa;
            font-size: 14px;
            margin-bottom: 16px;
        }
        .description {
            color: #ddd;
            line-height: 1.6;
            margin-bottom: 20px;
        }
        .actions {
            display: flex;
            gap: 12px;
            margin-bottom: 24px;
        }
        .btn {
            padding: 10px 24px;
            border-radius: 8px;
            border: none;
            cursor: pointer;
            font-size: 14px;
            font-weight: 500;
            transition: all 0.2s;
        }
        .btn-like {
            background: rgba(255,255,255,0.1);
            color: white;
        }
        .btn-like:hover {
            background: rgba(255,255,255,0.2);
        }
        .btn-like.liked {
            background: #e74c3c;
            color: white;
        }
        #ad-modal {
            display: none;
            position: fixed;
            top: 0;
            left: 0;
            width: 100%%;
            height: 100%%;
            background: rgba(0,0,0,0.95);
            z-index: 9999;
            align-items: center;
            justify-content: center;
        }
        #ad-modal.active {
            display: flex;
        }
        .ad-content {
            max-width: 800px;
            width: 90%%;
            text-align: center;
        }
        .ad-video {
            width: 100%%;
            border-radius: 12px;
        }
        .ad-info {
            margin-top: 16px;
            font-size: 14px;
            color: #aaa;
        }
        .skip-btn {
            margin-top: 16px;
            padding: 12px 32px;
            background: #fff;
            color: #000;
            border: none;
            border-radius: 8px;
            cursor: pointer;
            font-weight: 600;
            display: none;
        }
        .skip-btn.active {
            display: inline-block;
        }
    </style>
</head>
<body>
    <div class="container">
        <a href="/" class="back-link">← Back to Home</a>
        
        <div class="video-wrapper">
            <img id="thumbnail" src="%s" alt="Thumbnail" class="thumbnail">
            <div id="play-button" class="play-button"></div>
            <video id="video" controls controlsList="nodownload">
                <source src="/stream/%s" type="video/mp4">
            </video>
        </div>

        <div class="info">
            <h1 class="title">%s</h1>
            <div class="meta">
                👁️ %d views • ❤️ %d likes
            </div>
            <div class="actions">
                <button id="like-btn" class="btn btn-like">❤️ Like</button>
            </div>
            <div class="description">%s</div>
        </div>
    </div>

    <div id="ad-modal">
        <div class="ad-content">
            <video id="ad-video" class="ad-video" autoplay></video>
            <div class="ad-info">
                <span id="ad-countdown">Loading ad...</span>
            </div>
            <button id="skip-btn" class="skip-btn">Skip Ad</button>
        </div>
    </div>

    <script>
        const sessionId = localStorage.getItem('session_id') || crypto.randomUUID();
        localStorage.setItem('session_id', sessionId);

        const videoEl = document.getElementById('video');
        const thumbnail = document.getElementById('thumbnail');
        const playButton = document.getElementById('play-button');
        const adModal = document.getElementById('ad-modal');
        const adVideo = document.getElementById('ad-video');
        const adCountdown = document.getElementById('ad-countdown');
        const skipBtn = document.getElementById('skip-btn');
        const likeBtn = document.getElementById('like-btn');

        let watchStartTime = 0;
        let videoDuration = %d;

        // Play button click
        playButton.addEventListener('click', async () => {
            await showAd();
        });
        thumbnail.addEventListener('click', async () => {
            await showAd();
        });

        // Show ad
        async function showAd() {
            try {
                const adRes = await fetch('/api/ads/preroll');
                const adData = await adRes.json();
                
                if (adData.success && adData.data) {
                    const ad = adData.data;
                    adModal.classList.add('active');
                    adVideo.src = ad.content_url;
                    
                    let countdown = ad.duration_seconds || 7;
                    let canSkip = false;
                    
                    const timer = setInterval(() => {
                        countdown--;
                        if (countdown > 0) {
                            adCountdown.textContent = 'Ad will end in ' + countdown + 's';
                        } else {
                            clearInterval(timer);
                            skipBtn.classList.add('active');
                            adCountdown.textContent = 'You can skip now';
                            canSkip = true;
                        }
                    }, 1000);

                    skipBtn.addEventListener('click', () => {
                        if (canSkip) {
                            closeAd();
                            clearInterval(timer);
                        }
                    });

                    adVideo.addEventListener('ended', () => {
                        closeAd();
                        clearInterval(timer);
                    });

                    // Track ad impression
                    fetch('/api/ads/' + ad.id + '/impression', {
                        method: 'POST',
                        headers: {'Content-Type': 'application/json'},
                        body: JSON.stringify({
                            impression_type: 'view',
                            session_id: sessionId,
                            video_id: '%s'
                        })
                    });
                } else {
                    // No ad available, play video directly
                    closeAd();
                }
            } catch (err) {
                console.error('Ad error:', err);
                closeAd();
            }
        }

        function closeAd() {
            adModal.classList.remove('active');
            thumbnail.style.display = 'none';
            playButton.style.display = 'none';
            videoEl.classList.add('active');
            videoEl.play();
            watchStartTime = Date.now();
        }

        // Track watch duration
        videoEl.addEventListener('pause', trackView);
        videoEl.addEventListener('ended', trackView);
        window.addEventListener('beforeunload', trackView);

        function trackView() {
            if (watchStartTime === 0) return;
            const watchDuration = Math.floor((Date.now() - watchStartTime) / 1000);
            
            fetch('/api/videos/%s/view', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({
                    session_id: sessionId,
                    watch_duration: watchDuration,
                    video_duration: videoDuration
                })
            });
        }

        // Like button
        likeBtn.addEventListener('click', async () => {
            const token = localStorage.getItem('access_token');
            if (!token) {
                alert('Please login to like videos');
                return;
            }

            const method = likeBtn.classList.contains('liked') ? 'DELETE' : 'POST';
            const res = await fetch('/api/videos/%s/like', {
                method,
                headers: {'Authorization': 'Bearer ' + token}
            });

            if (res.ok) {
                likeBtn.classList.toggle('liked');
            }
        });

        // Check if already liked
        (async () => {
            const token = localStorage.getItem('access_token');
            if (token) {
                const res = await fetch('/api/videos/%s/liked', {
                    headers: {'Authorization': 'Bearer ' + token}
                });
                const data = await res.json();
                if (data.data && data.data.is_liked) {
                    likeBtn.classList.add('liked');
                }
            }
        })();
    </script>
</body>
</html>
	`, video.Title, video.ThumbnailURL, token, video.Title, video.ViewCount, video.LikeCount, video.Description, 
		video.DurationSeconds, video.ID.String(), video.ID.String(), video.ID.String(), video.ID.String())

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(html)
}

// ✅ FIXED: StreamVideo with context cancellation to prevent goroutine leak
func (h *StreamHandler) StreamVideo(c *fiber.Ctx) error {
	token := c.Params("token")

	// Get video by wrapper token
	video, err := h.videoService.GetVideoByWrapperToken(token)
	if err != nil {
		return c.Status(404).SendString("Video not found")
	}

	// ✅ Use request context for cancellation
	ctx := c.Context()

	// Create HTTP client with longer timeout for streaming
	client := &http.Client{
		Timeout: 10 * time.Minute, // ✅ Increased from 30s to 10 minutes
	}

	// ✅ Create request with context (enables cancellation)
	req, err := http.NewRequestWithContext(ctx, "GET", video.SourceURL, nil)
	if err != nil {
		return c.Status(500).SendString("Failed to create request")
	}

	// Forward Range header for seeking
	if rangeHeader := c.Get("Range"); rangeHeader != "" {
		req.Header.Set("Range", rangeHeader)
	}

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		// ✅ Check if context was cancelled
		if ctx.Err() == context.Canceled {
			return nil // User cancelled, clean exit
		}
		return c.Status(500).SendString("Failed to fetch video")
	}
	defer resp.Body.Close()

	// Set response headers
	c.Set("Content-Type", resp.Header.Get("Content-Type"))
	c.Set("Accept-Ranges", "bytes")
	c.Set("Cache-Control", "no-store")
	
	if contentLength := resp.Header.Get("Content-Length"); contentLength != "" {
		c.Set("Content-Length", contentLength)
	}
	
	if contentRange := resp.Header.Get("Content-Range"); contentRange != "" {
		c.Set("Content-Range", contentRange)
	}

	// Set status code
	c.Status(resp.StatusCode)

	// ✅ Stream with context awareness
	_, err = io.Copy(c.Response().BodyWriter(), resp.Body)
	if err != nil && ctx.Err() == context.Canceled {
		return nil // User cancelled, clean exit
	}
	
	return err
}