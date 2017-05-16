package com.mesosphere.sdk.api;

import org.eclipse.jetty.server.NetworkConnector;
import org.eclipse.jetty.server.Server;

import org.glassfish.jersey.jetty.JettyHttpContainerFactory;
import org.glassfish.jersey.server.ResourceConfig;

import javax.ws.rs.core.UriBuilder;
import java.net.URI;
import java.util.Collection;

/**
 * A JettyApiServer takes a list of POJO JAX-RS Resources and serves them on the indicated port of 0.0.0.0.
 */
public class JettyApiServer {
    private final Server server;
    private int port;

    public JettyApiServer(int port, Collection<Object> resources) {
        ResourceConfig resourceConfig = new ResourceConfig();
        for (Object resource : resources) {
            resourceConfig.register(resource);
        }

        URI baseUri = UriBuilder
                .fromUri("http://0.0.0.0/")
                .port(port).build();

        this.server = JettyHttpContainerFactory.createServer(baseUri, resourceConfig);
        if (port != 0) {
            this.port = port;
        }
    }

    public void start() throws Exception {
        try {
            server.start();
            server.dumpStdErr();
            server.join();
        } finally {
            server.destroy();
        }
    }

    public void stop() throws Exception {
        server.stop();
    }

    public boolean isStarted() {
        return server.isStarted();
    }

    public int getLocalPort() {
        if (this.port != 0) {
            return this.port;
        }
        NetworkConnector networkConnector = (NetworkConnector) server.getConnectors()[0];
        this.port = networkConnector.getLocalPort();
        return this.port;
    }
}
