package com.easyar.samples.cloud;

import org.asynchttpclient.AsyncCompletionHandler;
import org.asynchttpclient.AsyncHttpClient;
import org.asynchttpclient.DefaultAsyncHttpClient;
import org.asynchttpclient.Response;
import org.json.JSONObject;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.Base64;

/**
 * Created by qinsi on 6/12/16.
 */
public class AddTarget {

    private static final String HOST = "http://localhost:8888";
    private static final String APP_KEY = "test_app_key";
    private static final String APP_SECRET = "test_app_secret";

    public static void main(String[] args) throws IOException {
        AsyncHttpClient client = new DefaultAsyncHttpClient();

        JSONObject params = new JSONObject();
        params.put("image", Base64.getEncoder().encodeToString(
                Files.readAllBytes(Paths.get("test_target_image.jpg"))));
        params.put("foo", "bar");
        Auth.signParam(params, APP_KEY, APP_SECRET);

        client.preparePost(HOST + "/targets/")
                .setBody(params.toString().getBytes())
                .execute(new AsyncCompletionHandler<Void>() {
                    @Override
                    public Void onCompleted(Response response) throws Exception {
                        System.out.println(response.getResponseBody());
                        client.close();
                        return null;
                    }
                });
    }

}
