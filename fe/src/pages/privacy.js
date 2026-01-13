import React from 'react';
import Link from 'next/link';
import MainHeader from '@/components/MainHeader';
import Footer from '@/components/Footer';

export default function Privacy() {
  const contactEmail = process.env.NEXT_PUBLIC_CONTACT_EMAIL || 'contact@shorten.io';

  return (
    <div className="bg-[#1e1e1e] min-h-screen flex flex-col">
      <MainHeader />
      <div className="flex-1 px-6 py-20 lg:px-8">
        <div className="mx-auto max-w-4xl">
          <div className="text-center mb-12">
            <h1 className="text-[#f2bb4b] text-4xl font-semibold tracking-tight sm:text-5xl mb-4">
              Privacy Policy
            </h1>
            <p className="text-white text-sm sm:text-base">
              Last updated: {new Date().toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })}
            </p>
          </div>

          <div className="prose prose-invert max-w-none">
            <div className="bg-gray-800/50 rounded-lg p-6 sm:p-8 space-y-6 text-white">
              
              <section>
                <h2 className="text-[#f2bb4b] text-2xl font-semibold mb-4">1. Introduction</h2>
                <p className="text-gray-300 leading-relaxed">
                  Welcome to Shorten ("we," "our," or "us"). We are committed to protecting your privacy. 
                  This Privacy Policy explains how we collect, use, and safeguard information when you use 
                  our service to summarize web pages and generate text-to-speech audio.
                </p>
              </section>

              <section>
                <h2 className="text-[#f2bb4b] text-2xl font-semibold mb-4">2. Information We Collect</h2>
                <div className="space-y-4 text-gray-300">
                  <div>
                    <h3 className="text-white font-medium mb-2">2.1 Session Cookie</h3>
                    <p className="leading-relaxed">
                      We use a session cookie named "Shorten" to temporarily store the URL you provide 
                      for processing. This cookie is essential for the service to function and is automatically 
                      cleared after your session ends. The cookie is set with secure, httpOnly, and sameSite 
                      strict attributes for your security.
                    </p>
                  </div>
                  <div>
                    <h3 className="text-white font-medium mb-2">2.2 Analytics Data</h3>
                    <p className="leading-relaxed">
                      For analytics purposes, we temporarily store:
                    </p>
                    <ul className="list-disc list-inside ml-4 mt-2 space-y-1">
                      <li>URLs you request to summarize</li>
                      <li>Generated summaries</li>
                      <li>Text-to-speech audio files</li>
                      <li>IP addresses (for analytics only)</li>
                      <li>Performance metrics (latency data)</li>
                    </ul>
                  </div>
                  <div>
                    <h3 className="text-white font-medium mb-2">2.3 Google Analytics</h3>
                    <p className="leading-relaxed">
                      We use Google Analytics to understand how visitors interact with our website. 
                      Google Analytics collects anonymous usage data such as page views, time spent on pages, 
                      and navigation patterns. This data helps us improve our service.
                    </p>
                  </div>
                </div>
              </section>

              <section>
                <h2 className="text-[#f2bb4b] text-2xl font-semibold mb-4">3. How We Use Your Information</h2>
                <p className="text-gray-300 leading-relaxed mb-3">
                  We use the collected information solely for:
                </p>
                <ul className="list-disc list-inside ml-4 space-y-2 text-gray-300">
                  <li>Providing the summarization and text-to-speech services</li>
                  <li>Analyzing service performance and improving our algorithms</li>
                  <li>Understanding user behavior through Google Analytics</li>
                </ul>
              </section>

              <section>
                <h2 className="text-[#f2bb4b] text-2xl font-semibold mb-4">4. Data Retention</h2>
                <p className="text-gray-300 leading-relaxed">
                  All stored data (URLs, summaries, text-to-speech audio, and related analytics data) 
                  is automatically deleted after 7 days. We do not retain any personal information beyond 
                  this period.
                </p>
              </section>

              <section>
                <h2 className="text-[#f2bb4b] text-2xl font-semibold mb-4">5. Third-Party Services</h2>
                <div className="space-y-4 text-gray-300">
                  <div>
                    <h3 className="text-white font-medium mb-2">5.1 Google Analytics</h3>
                    <p className="leading-relaxed">
                      We use Google Analytics for website analytics. Google's use of data is governed by 
                      their Privacy Policy. You can opt-out of Google Analytics by installing the 
                      Google Analytics Opt-out Browser Add-on.
                    </p>
                  </div>
                  <div>
                    <h3 className="text-white font-medium mb-2">5.2 Groq (LLM Provider)</h3>
                    <p className="leading-relaxed">
                      We use Groq's services to generate summaries of web content. The URLs you provide 
                      are sent to Groq for processing. Please refer to Groq's Privacy Policy for information 
                      on how they handle data.
                    </p>
                  </div>
                </div>
              </section>

              <section>
                <h2 className="text-[#f2bb4b] text-2xl font-semibold mb-4">6. Data Security</h2>
                <p className="text-gray-300 leading-relaxed">
                  We implement appropriate technical and organizational measures to protect your data. 
                  However, no method of transmission over the Internet is 100% secure. While we strive 
                  to protect your information, we cannot guarantee absolute security.
                </p>
              </section>

              <section>
                <h2 className="text-[#f2bb4b] text-2xl font-semibold mb-4">7. No Personal Data Collection</h2>
                <p className="text-gray-300 leading-relaxed">
                  We do not require user registration or login. We do not collect personal information 
                  such as names, email addresses, or phone numbers. The service is completely free and 
                  anonymous.
                </p>
              </section>

              <section>
                <h2 className="text-[#f2bb4b] text-2xl font-semibold mb-4">8. Children's Privacy</h2>
                <p className="text-gray-300 leading-relaxed">
                  Our service is safe for everyone and does not knowingly collect personal information 
                  from children. Since we do not collect personal data, our service is accessible to 
                  users of all ages.
                </p>
              </section>

              <section>
                <h2 className="text-[#f2bb4b] text-2xl font-semibold mb-4">9. Changes to This Privacy Policy</h2>
                <p className="text-gray-300 leading-relaxed">
                  We may update this Privacy Policy from time to time. We will notify you of any changes 
                  by posting the new Privacy Policy on this page and updating the "Last updated" date. 
                  You are advised to review this Privacy Policy periodically for any changes.
                </p>
              </section>

            </div>
          </div>

          <div className="text-center mt-12">
            <Link 
              href="/"
              className="text-[#f2bb4b] hover:text-yellow-400 text-sm font-medium"
            >
              ← Back to Home
            </Link>
          </div>
        </div>
      </div>
      <Footer />
    </div>
  );
}