package com.easyar.samples.cloud;

import org.asynchttpclient.*;
import org.asynchttpclient.ws.WebSocket;
import org.asynchttpclient.ws.WebSocketTextListener;
import org.asynchttpclient.ws.WebSocketUpgradeHandler;
import org.json.JSONObject;
import org.msgpack.MessagePack;
import org.msgpack.packer.Packer;

import java.io.ByteArrayOutputStream;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.CountDownLatch;

/**
 * Created by qinsi on 6/12/16.
 */
public class SearchTarget {

    private static final String HOST = "localhost:8080";
    private static final String APP_KEY = "test_app_key";
    private static final String APP_SECRET = "test_app_secret";
    private static final int WEBSOCKET_FRAME_SIZE = 1024 * 1024;

    public static void main(String[] args) throws Exception {
        AsyncHttpClient client = new DefaultAsyncHttpClient(
                new DefaultAsyncHttpClientConfig.Builder()
                        .setWebSocketMaxFrameSize(WEBSOCKET_FRAME_SIZE)
                        .build());

        JSONObject params = new JSONObject();
        Auth.signParam(params, APP_KEY, APP_SECRET);

        JSONObject res = client.preparePost("http://" + HOST + "/tunnels/")
                .setBody(params.toString().getBytes())
                .execute(new AsyncCompletionHandler<JSONObject>() {
                    @Override
                    public JSONObject onCompleted(Response response) throws Exception {
                        return new JSONObject(response.getResponseBody());
                    }
                })
                .get();

        String host = res.getString("host");
        String port = res.getString("port");
        String tunnel = res.getJSONObject("result").getString("tunnel");

        Map<String, Object> map = new HashMap<>();
        byte[] image = Files.readAllBytes(Paths.get("test_search_image.jpg"));
        map.put("image", image);
        map.put("egg", "spam");

        MessagePack msgpack = new MessagePack();
        ByteArrayOutputStream out = new ByteArrayOutputStream();
        Packer packer = msgpack.createPacker(out);
        packer.write(map);
        byte[] data = out.toByteArray();
        out.close();

        CountDownLatch latch = new CountDownLatch(1);
        WebSocket ws = client.prepareGet(String.format("ws://%s:%s/services/recognize/%s", host, port, tunnel))
                .execute(new WebSocketUpgradeHandler.Builder().addWebSocketListener(
                        new WebSocketTextListener() {
                            @Override
                            public void onMessage(String s) {
                                System.out.println(s);
                                latch.countDown();
                            }

                            @Override
                            public void onOpen(WebSocket webSocket) {
                                webSocket.sendMessage(data);
                            }

                            @Override
                            public void onClose(WebSocket webSocket) {
                            }

                            @Override
                            public void onError(Throwable throwable) {
                            }
                        }
                ).build()).get();

        latch.await();
        client.close();
    }

}
